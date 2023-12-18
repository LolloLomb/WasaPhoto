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
	rt.router.GET("/user/:uid/stream", rt.wrap(rt.getMyStream))
	rt.router.DELETE("/user/:uid/following/:following_uid", rt.wrap(rt.unfollowUser))
	rt.router.POST("/user/:uid/ban", rt.wrap(rt.banUser))
	rt.router.DELETE("/user/:uid/ban/:banned_uid", rt.wrap(rt.unbanUser))
	rt.router.GET("/profile", rt.wrap(rt.getUserProfile))
	// VA TESTATO COME FUNZIONA, PER ORA È STATO TUTTO AGGIUNTO
	rt.router.POST("/photo", rt.wrap(rt.uploadPhoto))
	rt.router.POST("/photo/:photo_id/likes", rt.wrap(rt.likePhoto))
	rt.router.DELETE("/photo/:photo_id/likes/:uid", rt.wrap(rt.unlikePhoto))
	rt.router.POST("/photo/:photo_id/comment", rt.wrap(rt.commentPhoto))
	// QUANDO CANCELLERÒ IL COMMENTO AVRÒ SOLO IL COMMENT ID, FARE CONTROLLO CON AUTH PER VEDERE SE RICHIEDENTE E OWNER COINCIDONO
	rt.router.DELETE("/photo/:photo_id/comment/:comment_id", rt.wrap(rt.uncommentPhoto))
	rt.router.DELETE("/photo/:photo_id", rt.wrap(rt.deletePhoto))
	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
