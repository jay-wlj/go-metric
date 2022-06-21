package labels

import (
	"regexp"
	"strconv"
	"strings"
)

// See Reference: cn.lalaframework.jaf.monitor.core.util.TagReplaceUtil
// 实现一些通用的tags处理方法，在Http、sql埋点中使用，来过滤掉数字或查询语句

const (
	whiteSpace   = " "
	star         = "*"
	questionMark = "?"
)

var Filter = newFilter()

type filter struct {
	sqlPattern          *regexp.Regexp
	numberPattern       *regexp.Regexp
	whiteSpacesPattern  *regexp.Regexp
	singleQuotePattern  *regexp.Regexp
	timePattern         *regexp.Regexp
	questionMarkPattern *regexp.Regexp
}

func newFilter() *filter {
	return &filter{
		sqlPattern:          regexp.MustCompile("(in\\s*\\()(.*?)(\\))"),
		numberPattern:       regexp.MustCompile("\\d+"),
		whiteSpacesPattern:  regexp.MustCompile("[\\s]+"),
		singleQuotePattern:  regexp.MustCompile("((=|like)\\s*')(.*?)('\\s*)(,|$)"),
		timePattern:         regexp.MustCompile("\\?-\\?-\\?\\s\\?:\\?:\\?"),
		questionMarkPattern: regexp.MustCompile("(\\s*\\?\\s*,)+\\s*\\?"),
	}
}

func (f filter) FilterRoute(route string) string {
	return f.numberPattern.ReplaceAllString(route, star)
}

func (f filter) FilterResource(resource string) string {
	if resource == "" {
		return "-"
	}
	return resource
}

func (f filter) FilterHost(host string) string {
	// www.example.org:8080
	if strings.HasPrefix(host, "http://") {
		return strings.TrimLeft(host, "http://")
	} else if strings.HasPrefix(host, "https://") {
		return strings.TrimLeft(host, "https://")
	} else if host == "" {
		return "-"
	}
	return host
}

func (f filter) FilterStatusCode(statusCode int) string {
	if statusCode <= 0 || statusCode >= 600 {
		return "-"
	}
	return strconv.Itoa(statusCode)
}

func (f filter) FilterRet(ret string) string {
	_, err := strconv.ParseInt(strings.TrimSpace(ret), 10, 64)
	if err != nil {
		return "-"
	}
	return strings.TrimSpace(ret)
}

func (f filter) FilterSQL(sqlString string) (cmd string, sql string, ok bool) {
	// trim trailing space
	sqlString = strings.TrimSpace(strings.ToLower(sqlString))
	// lower case
	// replace multi white-spaces
	sqlString = f.whiteSpacesPattern.ReplaceAllString(sqlString, whiteSpace)
	// replace numbers
	sqlString = f.numberPattern.ReplaceAllString(sqlString, questionMark)
	// replace star
	sqlString = f.sqlPattern.ReplaceAllString(sqlString, "$1"+questionMark+"$3")
	//replace contents within single quotes
	sqlString = f.singleQuotePattern.ReplaceAllString(sqlString, "$1"+questionMark+"$4"+"$5")
	//replace multiple question marks
	sqlString = f.questionMarkPattern.ReplaceAllString(sqlString, questionMark)
	//replace time
	sqlString = f.timePattern.ReplaceAllString(sqlString, questionMark)

	sqlTypeAt := strings.Index(sqlString, whiteSpace)
	//
	if sqlTypeAt < 0 {
		return "", "", false
	}
	sqlType := sqlString[:sqlTypeAt]
	switch sqlType {
	case "insert":
		valuesAt := strings.Index(sqlString, "values")
		if valuesAt > 0 {
			sqlString = sqlString[:valuesAt-1] // truncate values
		}
		return sqlType, sqlString, len(sqlType) > 0 && len(sqlString) > 0
	case "select", "delete", "update":
		return sqlType, sqlString, len(sqlType) > 0 && len(sqlString) > 0
	default:
		return "", "", false
	}
}
