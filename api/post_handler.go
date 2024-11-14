package api

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	_ "net/http/httputil"
	"note_app/db"
	"note_app/fcm_notif"
	"note_app/types"
)

type PostHandler struct {
	postStore db.PostStore
	fcmClient *fcm_notif.FirebaseMessagingClient
	userStore db.UserStore
}

func NewPostHandler(postStore db.PostStore, userStore db.UserStore, fcmClient *fcm_notif.FirebaseMessagingClient) *PostHandler {
	return &PostHandler{
		postStore: postStore,
		fcmClient: fcmClient,
		userStore: userStore,
	}
}

// HandlePutPost UpdatePost update post
//
//	@Summary	Updating Post
//	@Tags		Posts
//	@Param		post	postID	path					types.PathParameter	true	"ID of post"
//	@Param		content	body	types.UpdatePostParams	true				"New content"
//	@Produce	json
//	@Success	200
//	@Failure	400	{string}	string
//	@Failure	404	{string}	string
//	@Router		/post/{id} [put]
func (h *PostHandler) HandlePutPost(c *fiber.Ctx) error {

	var (
		params types.UpdatePostParams
		postID = c.Params("id")
	)
	_, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return ErrBadRequest(err)
	}
	if err := c.BodyParser(&params); err != nil {
		return ErrBadRequest(err)
	}
	filter := db.Map{"_id": postID}
	if err := h.postStore.UpdatePost(c.Context(), filter, params); err != nil {
		return ErrNotResourceNotFound(err)
	}
	return c.JSON(map[string]string{"updated": postID})
}

// HandleDeletePost DeletePost Delete post
//
//	@Summary	Deleting Post
//	@Tags		Posts
//	@Param		post	postID	path	types.PathParameter	true	"ID of post"
//	@Produce	json
//	@Success	200	{map}		string
//	@Failure	400	{string}	string
//	@Failure	404	{string}	string
//	@Router		/post/{id} [delete]
func (h *PostHandler) HandleDeletePost(c *fiber.Ctx) error {
	postID := c.Params("id")
	_, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return ErrBadRequest(err)
	}
	if err := h.postStore.DeletePost(c.Context(), postID); err != nil {
		return ErrNotResourceNotFound(err)
	}
	return c.JSON(map[string]string{"deleted": postID})
}

// HandleInsertPost InsertPost Insert post
//
//	@Summary	Inserting Post
//	@Tags		Posts
//	@Param		post	body	types.CreatePostParams	true	"New Post"
//	@Produce	json
//	@Success	201	{object}	types.Post
//	@Failure	400	{string}	string
//	@Failure	500	{string}	string
//	@Router		/post [post]
func (h *PostHandler) HandleInsertPost(c *fiber.Ctx) error {
	var params types.CreatePostParams
	if err := c.BodyParser(&params); err != nil {
		return ErrBadRequest(err)
	}

	if errors := params.Validate(); len(errors) > 0 {
		return c.JSON(errors)
	}
	post := types.NewPostFromParams(params)
	insertedPost, err := h.postStore.InsertPost(c.Context(), post)
	if err != nil {
		return err
	}
	user, err := h.userStore.GetUser(c.Context(), params.Author)
	var friendsToken []string
	for _, friend := range user.Friends {
		userFriend, err := h.userStore.GetUserByObjectID(c.Context(), friend)
		if err != nil {
			return nil
		}
		friendsToken = append(friendsToken, userFriend.FCMToken)
	}

	err = h.fcmClient.SendNotification(c.Context(), friendsToken, "sa")
	if err != nil {
		return err
	}
	return c.JSON(insertedPost)
}

// HandleGetPost GetPost Get post
//
//	@Summary	Getting Post
//	@Tags		Posts
//	@Param		post	postID	path	types.PathParameter	true	"ID of post"
//	@Produce	json
//	@Success	200	{object}	types.Post
//	@Failure	400	{string}	string
//	@Failure	404	{string}	string
//	@Router		/post/{id} [get]
func (h *PostHandler) HandleGetPost(c *fiber.Ctx) error {
	var (
		postID = c.Params("id")
	)
	_, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return ErrBadRequest(err)
	}
	post, err := h.postStore.GetPostByID(c.Context(), postID)
	if err != nil {
		return ErrNotResourceNotFound(err)
	}
	return c.JSON(post)
}

// HandleGetPosts GetPosts Get posts
//
//	@Summary	Getting Posts
//	@Tags		Posts
//	@Produce	json
//	@Success	200	{array}		types.Post
//	@Failure	500	{string}	string
//	@Router		/posts [get]
func (h *PostHandler) HandleGetPosts(c *fiber.Ctx) error {
	posts, err := h.postStore.GetPosts(c.Context())
	if err != nil {
		return ErrNotResourceNotFound(err)
	}
	return c.JSON(posts)
}

// HandleGetPostsByUserID GetPosts get posts
//
//	@Summary	Getting posts from given user id
//	@Tags		Posts
//	@Param		user	userID	path	types.PathParameter	true	"ID of user"
//	@Produce	json
//	@Success	200	{array}		types.Post
//	@Failure	400	{string}	string
//	@Failure	404	{string}	string
//	@Router		/post/user/{id} [get]
func (h *PostHandler) HandleGetPostsByUserID(c *fiber.Ctx) error {
	var (
		userID = c.Params("id")
	)
	_, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return ErrBadRequest(err)
	}
	posts, err := h.postStore.GetPostsByUserID(c.Context(), userID)
	if err != nil {
		return ErrNotResourceNotFound(err)
	}
	return c.JSON(posts)
}
