package api

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"note_app/db"
	"note_app/types"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}

// HandlePutUser UpdateUser update user
//
//	@Summary	Updating user
//	@Tags		Users
//	@Param		post	userID	path					types.PathParameter	true	"ID of user"
//	@Param		content	body	types.UpdateUserParams	true				"User"
//	@Produce	json
//	@Success	200	{map}		string
//	@Failure	400	{string}	string
//	@Failure	404	{string}	string
//	@Router		/user/{id} [put]
func (h *UserHandler) HandlePutUser(c *fiber.Ctx) error {
	var (
		params types.UpdateUserParams
		userID = c.Params("id")
	)
	_, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return ErrBadRequest(err)
	}
	if err := c.BodyParser(&params); err != nil {
		return ErrBadRequest(err)
	}
	filter := db.Map{"_id": userID}
	if err := h.userStore.UpdateUser(c.Context(), filter, params); err != nil {
		return ErrNotResourceNotFound(err)
	}
	return c.JSON(map[string]string{"updated": userID})
}

// HandleDeleteUser DeleteUser Delete User
//
//	@Summary	Deleting User
//	@Tags		Users
//	@Param		user	userID	path	types.PathParameter	true	"ID of user"
//	@Produce	json
//	@Success	200	{map}		string
//	@Failure	400	{string}	string
//	@Failure	404	{string}	string
//	@Router		/user/{id} [delete]
func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	userID := c.Params("id")
	_, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return ErrBadRequest(err)
	}
	if err := h.userStore.DeleteUser(c.Context(), userID); err != nil {
		return err
	}
	return c.JSON(map[string]string{"deleted": userID})
}

// HandleInsertUser InsertUser Insert user
//
//	@Summary	Inserting user
//	@Tags		Users
//	@Param		user	body	types.CreateUserParams	true	"New user"
//	@Produce	json
//	@Success	201	{map}		string
//	@Failure	400	{string}	string
//	@Failure	500	{string}	string
//	@Router		/user [post]
func (h *UserHandler) HandleInsertUser(c *fiber.Ctx) error {
	var params types.CreateUserParams
	if err := c.BodyParser(&params); err != nil {
		return ErrBadRequest(err)
	}
	if errors := params.Validate(); len(errors) > 0 {
		return c.JSON(errors)
	}
	user, err := types.NewUserFromParams(params)
	if err != nil {
		return ErrBadRequest(err)
	}
	insertedUser, err := h.userStore.InsertUser(c.Context(), user)
	if err != nil {
		return err
	}
	return c.JSON(insertedUser)
}

// HandleGetUser GetUser Get user
//
//	@Summary	Getting user
//	@Tags		Users
//	@Param		user	userID	path	types.PathParameter	true	"ID of user"
//	@Produce	json
//	@Success	200	{map}		string
//	@Failure	400	{string}	string
//	@Failure	404	{string}	string
//	@Router		/user/{id} [get]
func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	var (
		userID = c.Params("id")
	)
	_, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return ErrBadRequest(err)
	}
	user, err := h.userStore.GetUser(c.Context(), userID)
	if err != nil {
		return ErrNotResourceNotFound(err)
	}
	return c.JSON(user)
}

// HandleGetUsers GetUsers Get users
//
//	@Summary	Getting users
//	@Tags		Users
//	@Produce	json
//	@Success	200	{map}		string
//	@Failure	404	{string}	string
//	@Router		/users [get]
func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	users, err := h.userStore.GetUsers(c.Context())
	if err != nil {
		return ErrNotResourceNotFound(err)
	}
	return c.JSON(users)
}

// HandleAddFriend AddFriend Add Friend
//
//	@Summary	Adding Freiend
//	@Tags		Users
//	@Param		post	userID	path				types.PathParameter	true	"ID of user"
//	@Param		userID	body	types.PathParameter	true				"User"
//	@Produce	json
//	@Success	200	{map}		string
//	@Failure	404	{string}	string
//	@Router		/user/{id}/add [put]
func (h *UserHandler) HandleAddFriend(c *fiber.Ctx) error {
	var (
		userID = c.Params("id")
		param  types.AddFriendParam
	)
	_, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return ErrBadRequest(err)
	}
	if err := c.BodyParser(&param); err != nil {
		return ErrBadRequest(err)
	}
	filter := db.Map{"_id": userID}
	err = h.userStore.AddFriend(c.Context(), filter, userID, param.UserID)
	if err != nil {
		return ErrNotResourceNotFound(err)
	}
	return c.JSON(map[string]string{"add friend": param.UserID})
}

// HandleRemoveFriend RemoveFriend Remove Friend
//
//	@Summary	Removing Freiend
//	@Tags		Users
//	@Param		post	userID	path				types.PathParameter	true	"ID of user"
//	@Param		userID	body	types.PathParameter	true				"User"
//	@Produce	json
//	@Success	200	{map}		string
//	@Failure	404	{string}	string
//	@Router		/user/{id}/remove [put]
func (h *UserHandler) HandleRemoveFriend(c *fiber.Ctx) error {
	var (
		userID = c.Params("id")
		param  types.AddFriendParam
	)
	_, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return ErrBadRequest(err)
	}
	if err := c.BodyParser(&param); err != nil {
		return ErrBadRequest(err)
	}
	filter := db.Map{"_id": userID}
	err = h.userStore.RemoveFriend(c.Context(), filter, userID, param.UserID)
	if err != nil {
		return ErrNotResourceNotFound(err)
	}
	return c.JSON(map[string]string{"remove friend": param.UserID})
}
