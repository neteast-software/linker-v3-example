// Package user 声明用户能力的 HTTP 入口。
//
// Route 是 HTTP 边界，不是单纯的转发层。每个 route 文件应集中保留路由声明、
// 资源范围、请求结构、请求头解析、响应结构和只服务于该 API 的短流程。
// 可复用的业务能力回到 user 根 package。
package user
