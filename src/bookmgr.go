package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"kkt.com/glog"
	"strings"
)

type contentSpecP struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type BooksInfo struct {
	Count       int              `json:"count"`
	Pages       int              `json:"pages"`
	Clazzs      []string         `json:"clazzs"`
	ContentSpec []ContentSpecCfg `json:"content_spec"`
	HotWords    []string         `json:"hot_words"`
}

type BookMgr struct {
	PageCount int
}

const defaultClazz = "default"

func NewBookMgr() (*BookMgr, error) {
	var mgr BookMgr
	mgr.PageCount = 20
	return &mgr, nil
}

func orderString(clazz string) (string, bool) {
	var orderMap = map[string]string{
		"recommend": "order by total_votes desc",
		"votes":     "order by total_votes desc",
		"reads":     "order by total_reads desc",
		"searches":  "order by total_searches desc",
		"chars":     "order by total_chars desc",
		"traces":    "order by total_reads desc",
		"score":     "order by score desc",
	}
	if _, ok := orderMap[clazz]; ok {
		return orderMap[clazz], true
	}
	return orderMap["score"], false
}

func (mgr *BookMgr) sqlRecommendString(clazz string, gender string, finished bool, curPage int) string {
	offset := curPage * mgr.PageCount
	orderStr, ok := orderString(clazz)
	sqlWhere := extraSqlWhereString(gender, finished)
	if !ok {
		if "" == sqlWhere {
			sqlWhere += "where "
		} else {
			sqlWhere += " and "
		}
		sqlWhere += "class='" + clazz + "'"
	}
	sqlExec := fmt.Sprintf(
		"select * from `books_table`%s %s limit %d offset %d",
		sqlWhere, orderStr, mgr.PageCount, offset)
	return sqlExec
}

func findClazzRecommendTableName(clazz string) string {
	var tableMap = map[string]string{
		"fprecommend":       "main_recommend_books",
		"girlrecommend":     "girl_recommend_books",
		"shelfrecommend":    "shelf_recommend_books",
		"directorrecommend": "director_recommend_books",
		"finishedrecommend": "finished_recommend_books",
	}
	return tableMap[clazz]
}

func (mgr *BookMgr) sqlClazzRecommendString(clazz string, gender string, finished bool) string {
	tableName := findClazzRecommendTableName(clazz)
	sqlWhere := extraSqlWhereString(gender, finished)
	sqlExec := fmt.Sprintf("select * from `books_table` a join `%s` b on a.`id`=b.`book_id`%s", tableName, sqlWhere)
	return sqlExec
}

func (mgr *BookMgr) sqlString(clazz string, gender string, finished bool, curPage int) string {
	switch clazz {
	case "fprecommend":
		fallthrough
	case "girlrecommend":
		fallthrough
	case "shelfrecommend":
		fallthrough
	case "directorrecommend":
		fallthrough
	case "finishedrecommend":
		return mgr.sqlClazzRecommendString(clazz, gender, finished)
	case "traces":
		fallthrough
	case "searches":
		fallthrough
	case "chars":
		fallthrough
	case "reads":
		fallthrough
	case "votes":
		fallthrough
	case "recommend":
		fallthrough
	default:
		return mgr.sqlRecommendString(clazz, gender, finished, curPage)
	}
	return mgr.sqlRecommendString("recommend", gender, finished, curPage)
}

func (mgr *BookMgr) queryCommonRecommendBooks(sqlExec string) ([]*Book, error) {
	var books = make([]*Book, 0)
	err := DBQuery(sqlExec, func(rows *sql.Rows) error {
		var book Book
		err := rows.Scan(&book.Id, &book.Name, &book.Abbreviation, &book.Author,
			&book.Cover, &book.AuthorAvatar, &book.Finished, &book.TotalReads,
			&book.TotalChars, &book.LastUpdateTime, &book.Class, &book.TotalSearches,
			&book.TotalVotes, &book.LastChapterTitle, &book.LastChapterUrl,
			&book.WithVIPChapter, &book.Gender, &book.Score, &book.Id)
		books = append(books, &book)
		return err
	})

	if nil != err {
		return books, err
	}

	return books, nil
}

func (mgr *BookMgr) queryDirectorRecommendBooks(sqlExec string) ([]*Book, error) {
	var books = make([]*Book, 0)
	err := DBQuery(sqlExec, func(rows *sql.Rows) error {
		var book Book
		err := rows.Scan(&book.Id, &book.Name, &book.Abbreviation, &book.Author,
			&book.Cover, &book.AuthorAvatar, &book.Finished, &book.TotalReads,
			&book.TotalChars, &book.LastUpdateTime, &book.Class, &book.TotalSearches,
			&book.TotalVotes, &book.LastChapterTitle, &book.LastChapterUrl,
			&book.WithVIPChapter, &book.Gender, &book.Score, &book.Id,
			&book.RWords, &book.RUser)
		books = append(books, &book)
		return err
	})

	if nil != err {
		return books, err
	}

	return books, nil
}

func (mgr *BookMgr) queryBooks(sqlExec string) ([]*Book, error) {
	var books = make([]*Book, 0)
	err := DBQuery(sqlExec, func(rows *sql.Rows) error {
		var book Book
		err := rows.Scan(&book.Id, &book.Name, &book.Abbreviation, &book.Author,
			&book.Cover, &book.AuthorAvatar, &book.Finished, &book.TotalReads,
			&book.TotalChars, &book.LastUpdateTime, &book.Class, &book.TotalSearches,
			&book.TotalVotes, &book.LastChapterTitle, &book.LastChapterUrl,
			&book.WithVIPChapter, &book.Gender, &book.Score)
		books = append(books, &book)
		return err
	})

	if nil != err {
		return books, err
	}

	return books, nil
}

func (mgr *BookMgr) queryBooksList(clazz string, sqlExec string) ([]*Book, error) {
	switch clazz {
	case "fprecommend":
		fallthrough
	case "girlrecommend":
		fallthrough
	case "shelfrecommend":
		fallthrough
	case "finishedrecommend":
		return mgr.queryCommonRecommendBooks(sqlExec)
	case "directorrecommend":
		return mgr.queryDirectorRecommendBooks(sqlExec)
	case "traces":
		fallthrough
	case "searches":
		fallthrough
	case "chars":
		fallthrough
	case "reads":
		fallthrough
	case "votes":
		fallthrough
	case "recommend":
		fallthrough
	default:
		return mgr.queryBooks(sqlExec)
	}
	return mgr.queryBooks(sqlExec)
}

func extraSqlWhereString(gender string, finished bool) string {
	sqlStr := ""
	if gender == "girl" || finished {
		sqlStr += " where "
	}
	if gender == "girl" {
		sqlStr += "gender='girl'"
	}
	if gender == "girl" && finished {
		sqlStr += " and "
	}
	if finished {
		sqlStr += "finished=1"
	}
	return sqlStr
}

func (mgr *BookMgr) QueryBooksList(clazz string, gender string, finished bool, curPage int) ([]*Book, error) {
	if curPage < 0 {
		curPage = 0
	}
	sqlExec := mgr.sqlString(clazz, gender, finished, curPage)

	return mgr.queryBooksList(clazz, sqlExec)
}

func (mgr *BookMgr) QueryBooksInfo(clazz string, gender string, finished bool) (*BooksInfo, error) {
	sqlExec := "select count(*) from `books_table`"
	sqlWhere := extraSqlWhereString(gender, finished)
	sqlExec += sqlWhere

	var info BooksInfo
	err := DBQuery(sqlExec, func(rows *sql.Rows) error {
		err := rows.Scan(&info.Count)
		return err
	})

	if nil != err {
		return nil, err
	}

	sqlExec = "select distinct class from `books_table`" + sqlWhere
	info.Clazzs = make([]string, 0)

	err = DBQuery(sqlExec, func(rows *sql.Rows) error {
		bytes := make([]byte, 128)
		err := rows.Scan(&bytes)
		if nil == err && 0 != len(bytes) {
			info.Clazzs = append(info.Clazzs, string(bytes))
		}
		return err
	})

	if nil != err {
		return nil, err
	}

	sqlExec = fmt.Sprintf("select word from `search_words_table` order by count limit 20")
	info.HotWords = make([]string, 0)
	err = DBQuery(sqlExec, func(rows *sql.Rows) error {
		word := ""
		err := rows.Scan(&word)
		if nil == err {
			info.HotWords = append(info.HotWords, word)
		}
		return err
	})

	if 0 == len(info.HotWords) {
		info.HotWords = []string{"唐家三少", "修真", "武侠仙侠"}
	}
	info.Pages = info.Count/mgr.PageCount + 1
	info.ContentSpec = cfg.ContentSpec

	return &info, nil
}

type recommendBook struct {
	Id     string `json:"id"`
	RWords string `json:"rwords"`
	RUser  string `json:"ruser"`
}

type recommendP struct {
	Clazz string          `json:"clazz"`
	Books []recommendBook `json:"books"`
}

func sqlFromRecommendParameter(p recommendP, director bool) (string, error) {
	var valueStr = ""
	var count = 0

	for _, book := range p.Books {
		if 0 < count {
			valueStr += ","
		}
		count += 1
		valueStr += "("
		valueStr += fmt.Sprintf("'%s'", book.Id)
		if director {
			if "" == book.RWords || "" == book.RUser {
				return "", errors.New("Invalid parameter")
			}
			valueStr += fmt.Sprintf(",'%s'", book.RWords)
			valueStr += fmt.Sprintf(",'%s'", book.RUser)
		}
		valueStr += ")"
	}

	var sqlExec = ""
	tableName := findClazzRecommendTableName(p.Clazz)
	if director {
		sqlExec = fmt.Sprintf("insert into `%s` (book_id, rwords, ruser) values %s", tableName, valueStr)
	} else {
		sqlExec = fmt.Sprintf("insert into `%s` (book_id) values %s", tableName, valueStr)
	}

	return sqlExec, nil
}

func (mgr *BookMgr) sqlClazzRecommendSetString(p recommendP) (string, error) {
	switch p.Clazz {
	case "directorrecommend":
		return sqlFromRecommendParameter(p, true)
	default:
		return sqlFromRecommendParameter(p, false)
	}
}

func (mgr *BookMgr) SetBooks(key string, body interface{}) error {
	jsonStr, err := json.Marshal(body)
	if nil != err {
		return err
	}

	var p recommendP
	err = json.Unmarshal(jsonStr, &p)
	if nil != err {
		return err
	}

	tableName := findClazzRecommendTableName(p.Clazz)
	sqlReset := fmt.Sprintf("truncate table %s", tableName)
	err = DBExec(sqlReset)
	if nil != err {
		return err
	}

	sqlInsert, err := mgr.sqlClazzRecommendSetString(p)
	if nil != err {
		return err
	}
	return DBExec(sqlInsert)
}

func (mgr *BookMgr) GetBookChapters(name string, author string) ([]*Chapter, error) {
	tableName := sha256.Sum256([]byte(name + author))
	sqlExec := fmt.Sprintf("select * from `%s`", hex.EncodeToString(tableName[0:]))

	var chapters = make([]*Chapter, 0)
	err := DBQuery(sqlExec, func(rows *sql.Rows) error {
		var chapter Chapter
		err := rows.Scan(&chapter.NativeId, &chapter.Id, &chapter.Title, &chapter.Url, &chapter.Vip)
		chapters = append(chapters, &chapter)
		return err
	})
	if nil != err {
		return make([]*Chapter, 0), err
	}
	return chapters, nil
}

type searchKeyP struct {
	word  string
	count int
}

func (mgr *BookMgr) updateSearches(books []*Book, key string) {
	sqlExec := fmt.Sprintf("update `books_table` set total_searches = case id ")
	ids := ""
	for i, book := range books {
		book.TotalSearches += 1
		sqlExec += fmt.Sprintf("when '%s' then '%d' ", book.Id, book.TotalSearches)
		if 0 < i {
			ids += ","
		}
		ids += fmt.Sprintf("'%s'", book.Id)
	}
	sqlExec += fmt.Sprintf("end where id in (%s)", ids)
	err := DBExec(sqlExec)
	if nil != err {
		glog.Error("Error: fail to update searches")
	}

	p := searchKeyP{word: key, count: 0}
	sqlExec = fmt.Sprintf("select * from `search_words_table` where word='%s'", key)
	err = DBQuery(sqlExec, func(rows *sql.Rows) error {
		err := rows.Scan(&p.word, &p.count)
		return err
	})

	p.count += 1
	sqlExec = fmt.Sprintf("insert into `search_words_table` values ('%s', %d)"+
		" on duplicate key update count=%d", p.word, p.count, p.count)
	err = DBExec(sqlExec)
	if nil != err {
		glog.Error("Error: fail to update search words")
	}
}

func (mgr *BookMgr) SearchBooks(clazz string, curPage int) ([]*Book, int, error) {
	sqlWhere := ""
	key := ""
	cs := strings.Split(clazz, "")
	for _, c := range cs {
		key += fmt.Sprintf("%%%s", c)
	}
	if "" != clazz {
		key += fmt.Sprintf("%%")
		sqlWhere = " where author like '" + key +
			"' or name like '" + key +
			"' or class like '" + key + "'"
	}
	if curPage < 0 {
		curPage = 0
	}

	count := -1
	if 0 == curPage {
		sqlExec := fmt.Sprintf("select count(*) from `books_table`%s order by score desc", sqlWhere)
		err := DBQuery(sqlExec, func(rows *sql.Rows) error {
			err := rows.Scan(&count)
			return err
		})
		if nil != err {
			return make([]*Book, 0), -1, err
		}
	}

	offset := curPage * mgr.PageCount
	sqlExec := fmt.Sprintf("select * from `books_table`%s order by score desc limit %d offset %d",
		sqlWhere, mgr.PageCount, offset)
	books, err := mgr.queryBooks(sqlExec)
	if nil != err {
		return make([]*Book, 0), count, err
	}

	go mgr.updateSearches(books, clazz)

	return books, count, nil
}

type bookReqBodyBaseP struct {
	BookId   string `json:"book_id"`
	ClientId string `json:"client_id"`
}

func (p *bookReqBodyBaseP) validate() bool {
	if "" == p.BookId {
		return false
	}
	if "" == p.ClientId {
		return false
	}
	return true
}

func (mgr *BookMgr) addBookRead(body interface{}) (*Book, error) {
	bodyJSON, err := json.Marshal(body)
	if nil != err {
		return nil, err
	}
	var p bookReqBodyBaseP
	err = json.Unmarshal(bodyJSON, &p)
	if nil != err {
		return nil, err
	}
	if !p.validate() {
		return nil, errors.New("Invalid parameter")
	}

	sqlExec := fmt.Sprintf("select * from `books_table` where id='%s'", p.BookId)
	books, err := mgr.queryBooks(sqlExec)
	if nil != err {
		return nil, err
	}
	if 0 == len(books) {
		return nil, errors.New("Invalid parameter")
	}

	books[0].TotalReads += 1
	sqlExec = fmt.Sprintf("update `books_table` set total_reads='%d' where id='%s'",
		books[0].TotalReads, p.BookId)
	err = DBExec(sqlExec)
	if nil != err {
		return nil, err
	}

	return books[0], nil
}
