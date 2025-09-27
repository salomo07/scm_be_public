package consts

const QueryInit = `{"selector":{"type":"initialdata"}}`

func BodySecurity(dbName string) string {
	return `{"admins": {"names": ["` + dbName + `"],"roles": ["admin_role"]},"members": {"names": [],"roles": []}}`
}
func QueryCompanyAlias(alias string, appid string) string {
	return `{"selector": {"table":"company","alias":"` + alias + `","appid":"` + appid + `"},"limit":1}`
}
func QueryFindCompany(idCompany string) string {
	return `{"_id": "` + idCompany + `"}`
}
