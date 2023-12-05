package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/", rt.getHelloWorld)
	rt.router.GET("/context", rt.wrap(rt.getContextReply))

	rt.router.POST("/session", rt.login)
	rt.router.PUT("/user/:uid/username", rt.wrap(rt.setMyUserName))
	rt.router.POST("/user/:uid/following", rt.wrap(rt.followUser))
	//rt.router.GET("/user/:uid/stream", rt.wrap(rt.getMyStream))
	rt.router.DELETE("/user/:uid/following/:following_uid", rt.wrap(rt.unfollowUser))
	rt.router.POST("/user/:uid/ban", rt.wrap(rt.banUser))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
