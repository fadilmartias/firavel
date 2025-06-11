package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8" // atau v9
	"gorm.io/gorm"
)

// QueryParams mem-parsing dan menyimpan parameter dari URL.
type QueryParams struct {
	Fields  []string
	Joins   []string
	Orders  []string
	Groups  []string
	Filters url.Values // Menggunakan url.Values untuk menangani format seperti filters[field][op]=value
	Limit   int
	Page    int
}

// NewQueryParams membuat instance QueryParams dari URL query.
func NewQueryParams(query url.Values) *QueryParams {
	// Fungsi helper untuk mengambil nilai atau default string kosong
	get := func(key string) string {
		return query.Get(key)
	}

	// Parsing limit dan page
	limit, _ := strconv.Atoi(get("limit"))
	if limit == 0 && get("limit") != "0" {
		limit = 10 // Default limit
	}

	page, _ := strconv.Atoi(get("page"))
	if page == 0 {
		page = 1 // Default page
	}

	// Parsing filter
	filters := url.Values{}
	for key, values := range query {
		if strings.HasPrefix(key, "filters[") {
			filters[key] = values
		}
	}

	// Fungsi helper untuk memisahkan string dengan koma
	split := func(s string) []string {
		if s == "" {
			return nil
		}
		return strings.Split(s, ",")
	}

	return &QueryParams{
		Fields:  split(get("fields")),
		Joins:   split(get("joins")),
		Orders:  split(get("orders")),
		Groups:  split(get("groups")),
		Filters: filters,
		Limit:   limit,
		Page:    page,
	}
}

/**
 * Membangun kueri GORM secara dinamis dari parameter query request.
 * Mendukung: fields (Select), joins (Preload), orders (Order), filters (Where & Joins), groups (Group), dan pagination (Limit & Offset).
 *
 * @param {*gorm.DB} db - Instance GORM awal yang akan dimodifikasi.
 * @param {url.Values} query - Nilai-nilai dari URL query (misalnya, c.Request.URL.Query()).
 * @param {bool} isSingle - Jika true, limit & offset tidak akan diterapkan (untuk First, Take).
 * @returns {*gorm.DB} Instance GORM yang telah dimodifikasi dan siap dieksekusi.
 */
func BuildGormQuery(db *gorm.DB, query url.Values, isSingle bool) *gorm.DB {
	params := NewQueryParams(query)

	// Terapkan semua builder secara berurutan
	db = applyFilters(db, params.Filters)
	db = applySelectsAndPreloads(db, params.Fields, params.Joins)
	db = applyOrders(db, params.Orders)
	db = applyGroups(db, params.Groups)
	db = applyPagination(db, params.Limit, params.Page, isSingle)

	return db
}

// applyFilters menerapkan klausa WHERE dan JOINs berdasarkan parameter filter.
func applyFilters(db *gorm.DB, filters url.Values) *gorm.DB {
	// Regex untuk mengekstrak field dan operator dari kunci filter, contoh: filters[users.name][like]
	filterRegex := regexp.MustCompile(`^filters\[([^\]]+)\](?:\[([^\]]+)\])?$`)

	// Peta untuk melacak join yang sudah ditambahkan agar tidak duplikat
	addedJoins := make(map[string]bool)

	for key, values := range filters {
		matches := filterRegex.FindStringSubmatch(key)
		if len(matches) < 2 {
			continue
		}

		field := matches[1]
		// Operator default adalah 'eq' jika tidak dispesifikkan
		operator := "eq"
		if len(matches) > 2 && matches[2] != "" {
			operator = matches[2]
		}

		value := values[0]

		// Periksa apakah ini filter untuk relasi (misal: "users.name")
		if strings.Contains(field, ".") {
			parts := strings.SplitN(field, ".", 2)
			relationName := strings.Title(parts[0]) // Konversi "users" menjadi "User" untuk nama relasi

			// Tambahkan JOIN jika belum ada
			if _, exists := addedJoins[relationName]; !exists {
				// Penting: "Joins" di GORM menggunakan nama field struct, bukan nama tabel.
				db = db.Joins(relationName)
				addedJoins[relationName] = true
			}
		}

		// Bangun kondisi WHERE
		// Peta untuk mengubah operator dari URL ke SQL
		opMap := map[string]string{
			"eq": "=", "neq": "!=", "gt": ">", "gte": ">=",
			"lt": "<", "lte": "<=", "like": "LIKE", "ilike": "ILIKE",
		}

		sqlOp, isValidOp := opMap[strings.ToLower(operator)]
		if !isValidOp {
			continue // Abaikan operator yang tidak didukung
		}

		// Handle nilai 'null'
		var queryValue interface{}
		if strings.ToLower(value) == "null" {
			if sqlOp == "=" {
				db = db.Where(fmt.Sprintf("%s IS NULL", field))
			} else if sqlOp == "!=" {
				db = db.Where(fmt.Sprintf("%s IS NOT NULL", field))
			}
			continue
		} else {
			queryValue = value
		}

		// Untuk operator LIKE, tambahkan wildcard '%'
		if sqlOp == "LIKE" || sqlOp == "ILIKE" {
			queryValue = "%" + value + "%"
		}

		// Terapkan kondisi Where
		db = db.Where(fmt.Sprintf("%s %s ?", field, sqlOp), queryValue)
	}
	return db
}

// applySelectsAndPreloads menerapkan Preload (untuk joins) dan Select (untuk fields).
func applySelectsAndPreloads(db *gorm.DB, fields []string, joins []string) *gorm.DB {
	// --- Penanganan Fields (Select) ---
	mainModelFields := []string{}
	preloadFields := make(map[string][]string)

	if len(fields) > 0 {
		hasSpecificField := false
		for _, field := range fields {
			if strings.Contains(field, ".") {
				parts := strings.SplitN(field, ".", 2)
				relationName := strings.Title(parts[0]) // users.name -> User
				fieldName := parts[1]
				preloadFields[relationName] = append(preloadFields[relationName], fieldName)
			} else {
				mainModelFields = append(mainModelFields, field)
				if field != "id" {
					hasSpecificField = true
				}
			}
		}

		// Selalu sertakan 'id' jika field lain dari model utama diminta
		if hasSpecificField && !contains(mainModelFields, "id") {
			mainModelFields = append([]string{"id"}, mainModelFields...)
		}

		if len(mainModelFields) > 0 {
			db = db.Select(mainModelFields)
		}
	}

	// --- Penanganan Joins (Preload) ---
	for _, join := range joins {
		// GORM menangani nested preload dengan format "Relation1.Relation2"
		relationPath := ""
		for _, part := range strings.Split(join, ".") {
			relationName := strings.Title(part)
			if relationPath != "" {
				relationPath += "."
			}
			relationPath += relationName

			// Jika ada fields spesifik untuk preload ini, terapkan dengan custom scope
			if fields, ok := preloadFields[relationName]; ok {
				// Selalu sertakan 'id' di relasi jika field lain diminta
				if !contains(fields, "id") {
					fields = append([]string{"id"}, fields...)
				}

				db = db.Preload(relationPath, func(db *gorm.DB) *gorm.DB {
					return db.Select(fields)
				})
			} else {
				// Jika tidak ada fields spesifik, preload semua kolom
				db = db.Preload(relationPath)
			}
		}
	}

	return db
}

// applyOrders menerapkan klausa ORDER BY.
func applyOrders(db *gorm.DB, orders []string) *gorm.DB {
	for _, orderItem := range orders {
		parts := strings.Split(orderItem, ":")
		field := parts[0]
		direction := "ASC"
		if len(parts) > 1 && strings.ToLower(parts[1]) == "desc" {
			direction = "DESC"
		}
		// Penting: Sama seperti filter, order pada relasi butuh nama tabel.
		// Contoh: orders=users.name:desc
		db = db.Order(fmt.Sprintf("%s %s", field, direction))
	}
	return db
}

// applyGroups menerapkan klausa GROUP BY.
func applyGroups(db *gorm.DB, groups []string) *gorm.DB {
	if len(groups) > 0 {
		db = db.Group(strings.Join(groups, ", "))
	}
	return db
}

// applyPagination menerapkan klausa LIMIT dan OFFSET.
func applyPagination(db *gorm.DB, limit, page int, isSingle bool) *gorm.DB {
	if isSingle || limit == 0 {
		return db
	}

	offset := (page - 1) * limit
	return db.Limit(limit).Offset(offset)
}

// contains adalah fungsi helper untuk memeriksa keberadaan string dalam slice.
func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}

// Pagination berisi detail paginasi untuk respons.
type Pagination struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int64 `json:"total_pages"`
	TotalItems int64 `json:"total_items"`
}

// PaginatedResponse adalah struktur untuk data dengan paginasi (isSingle = false).
type PaginatedResponse[T any] struct {
	Pagination Pagination `json:"pagination"`
	Data       []T        `json:"data"`
}

// SingleResponse adalah struktur untuk data tunggal (isSingle = true).
type SingleResponse[T any] struct {
	Data T `json:"data"`
}

/**
 * Mengambil data dari database dengan logika caching Redis.
 * Fungsi ini generik dan dapat bekerja dengan model GORM apa pun.
 *
 * @template T - Tipe struct model GORM (misal: User, Post).
 * @param {*redis.Client} redisClient - Instance klien Redis yang aktif.
 * @param {*gorm.DB} db - Instance kueri GORM yang sudah dibangun (oleh BuildGormQuery).
 * @param {*QueryParams} params - Parameter query yang sudah diparsing, diperlukan untuk paginasi.
 * @param {string} cacheKey - Kunci unik untuk caching di Redis. Jika string kosong, caching dilewati.
 * @param {time.Duration} cacheDuration - Durasi cache di Redis (misal: 60 * time.Second).
 * @param {bool} isSingle - Jika true, akan mengambil satu record (First); jika false, mengambil slice dengan paginasi.
 * @returns {any, error} - Mengembalikan PaginatedResponse[T] atau SingleResponse[T] dalam bentuk 'any', dan error jika terjadi.
 */
func FetchAndCache[T any](
	ctx context.Context, // Menerima context sebagai parameter pertama
	redisClient *redis.Client,
	db *gorm.DB,
	params *QueryParams,
	cacheKey string,
	cacheDuration time.Duration,
	isSingle bool,
) (any, error) {

	// 1. Coba ambil dari Cache menggunakan context yang diberikan
	if cacheKey != "" {
		// Gunakan 'ctx' dari parameter
		cachedData, err := redisClient.Get(ctx, cacheKey).Result()
		if err == nil {
			// ... (logika unmarshal cache hit tidak berubah)
			if isSingle {
				var response SingleResponse[T]
				if json.Unmarshal([]byte(cachedData), &response) == nil {
					return response, nil
				}
			} else {
				var response PaginatedResponse[T]
				if json.Unmarshal([]byte(cachedData), &response) == nil {
					return response, nil
				}
			}
		}
	}

	// 2. Jika Cache Miss, ambil dari Database
	var response any
	var dbErr error

	// Penting: Terapkan context ke GORM
	// Ini akan membuat kueri GORM dapat dibatalkan jika context-nya selesai.
	db = db.WithContext(ctx)

	if isSingle {
		var result T
		dbErr = db.First(&result).Error // Kueri ini sekarang cancellable
		if dbErr != nil {
			return nil, dbErr
		}
		response = SingleResponse[T]{Data: result}
	} else {
		// ... (logika paginasi tidak berubah)
		var items []T
		var totalItems int64
		model := new(T)

		dbErr = db.Model(model).Count(&totalItems).Error
		if dbErr != nil {
			return nil, dbErr
		}

		paginatedDb := applyPagination(db, params.Limit, params.Page, false)
		dbErr = paginatedDb.Find(&items).Error
		if dbErr != nil {
			return nil, dbErr
		}

		// ... (logika pembuatan objek paginasi tidak berubah)
		totalPages := int64(0)
		if params.Limit > 0 {
			totalPages = int64(math.Ceil(float64(totalItems) / float64(params.Limit)))
		}
		pagination := Pagination{
			Page:       params.Page,
			PageSize:   params.Limit,
			TotalPages: totalPages,
			TotalItems: totalItems,
		}
		response = PaginatedResponse[T]{
			Pagination: pagination,
			Data:       items,
		}
	}

	// 3. Simpan hasil ke Cache jika perlu
	if cacheKey != "" && dbErr == nil {
		jsonResponse, err := json.Marshal(response)
		if err == nil {
			// Gunakan 'ctx' dari parameter
			redisClient.SetEX(ctx, cacheKey, jsonResponse, cacheDuration)
		}
	}

	return response, nil
}
