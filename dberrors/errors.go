package dberrors

import "errors"

const UniqueConstraint = "23505"
const NotNullConstraint = "23502"

var ErrUserNotFound = errors.New("UserNotFound")
var ErrUserExists = errors.New("UserExists")
var ErrUserConflict = errors.New("UserConflict")

var ErrForumNotFound = errors.New("ForumNotFound")
var ErrForumExists = errors.New("ForumExists")

var ErrThreadExists = errors.New("ThreadExists")
var ErrThreadNotFound = errors.New("ErrThreadNotFound")

var ErrPostsConflict = errors.New("ErrPostsConflict")
