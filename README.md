## Backend service for [splatapp](https://github.com/hauchongtang/studytgt)

## Testing the API
Head over to [swagger](https://splatbackend-production.up.railway.app/docs/index.html#/default)

## Version: 1.0

**License:** [MIT](https://opensource.org/licenses/MIT)

### /cached/users

#### GET
##### Summary

Get all users from cache

##### Description

Gets all users from cache.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| token | header | Authorization token | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [ [controllers.userType](#controllersusertype) ] |
| 404 | Not Found | [controllers.errorResult](#controllerserrorresult) |

##### Security

| Security Schema | Scopes |
| --- | --- |
| ApiKeyAuth | |

### /cached/users/{id}

#### GET
##### Summary

Get a User by id from cache

##### Description

Gets a user from the cache if there is a hit. This is the default endpoint.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path | userId | Yes | string |
| token | header | Authorization token | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [controllers.userType](#controllersusertype) |
| 404 | Not Found | [controllers.errorResult](#controllerserrorresult) |

##### Security

| Security Schema | Scopes |
| --- | --- |
| ApiKeyAuth | |

### /tasks

#### GET
##### Summary

Get all task activities

##### Description

Gets all tasks from the database. Represents all activities.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| token | header | Authorization token | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [ [controllers.taskType](#controllerstasktype) ] |
| 404 | Not Found | [controllers.errorResult](#controllerserrorresult) |

##### Security

| Security Schema | Scopes |
| --- | --- |
| ApiKeyAuth | |

#### POST
##### Summary

Add a task

##### Description

Adds task to the database

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| data | body | Task details | Yes | [controllers.taskAddType](#controllerstaskaddtype) |
| token | header | Authorization token | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [controllers.taskType](#controllerstasktype) |
| 400 | Bad Request | [controllers.errorResult](#controllerserrorresult) |

##### Security

| Security Schema | Scopes |
| --- | --- |
| ApiKeyAuth | |

### /tasks/{id}

#### GET
##### Summary

Get all Tasks of a particular user

##### Description

Gets tasks of a particular user via userId.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path | taskId | Yes | string |
| token | header | Authorization token | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [ [controllers.taskType](#controllerstasktype) ] |
| 404 | Not Found | [controllers.errorResult](#controllerserrorresult) |

##### Security

| Security Schema | Scopes |
| --- | --- |
| ApiKeyAuth | |

#### PUT
##### Summary

Sets a task to be hidden

##### Description

Updates the task via provided taskId to be hidden.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path | userId | Yes | string |
| token | header | Authorization token | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [controllers.taskType](#controllerstasktype) |
| 404 | Not Found | [controllers.errorResult](#controllerserrorresult) |

##### Security

| Security Schema | Scopes |
| --- | --- |
| ApiKeyAuth | |

### /users

#### GET
##### Summary

Get all users

##### Description

Gets all users from database directly. Use it to test whether cache is updated correctly.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| token | header | Authorization token | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [ [controllers.userType](#controllersusertype) ] |
| 404 | Not Found | [controllers.errorResult](#controllerserrorresult) |

##### Security

| Security Schema | Scopes |
| --- | --- |
| ApiKeyAuth | |

### /users/{id}

#### DELETE
##### Summary

Delete a user given a userId

##### Description

Deletes a user via userId. Only admin access.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path | userId | Yes | string |
| adminId | query | adminId | Yes | string |
| token | header | Authorization token | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [controllers.userType](#controllersusertype) |
| 404 | Not Found | [controllers.errorResult](#controllerserrorresult) |

##### Security

| Security Schema | Scopes |
| --- | --- |
| ApiKeyAuth | |

#### GET
##### Summary

Get a User by id from database

##### Description

Gets a user from database. Use this to check if the cache is updated compared to the database.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path | userId | Yes | string |
| token | header | Authorization token | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [controllers.userType](#controllersusertype) |
| 404 | Not Found | [controllers.errorResult](#controllerserrorresult) |

##### Security

| Security Schema | Scopes |
| --- | --- |
| ApiKeyAuth | |

#### PUT
##### Summary

Increase points of a user

##### Description

Increase points of a user by specified amount.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path | userId | Yes | string |
| pointstoadd | query | pointsToAdd | Yes | string |
| token | header | Authorization token | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [controllers.userType](#controllersusertype) |
| 404 | Not Found | [controllers.errorResult](#controllerserrorresult) |

##### Security

| Security Schema | Scopes |
| --- | --- |
| ApiKeyAuth | |

### /users/login

#### POST
##### Summary

User log in

##### Description

Responds with user details, including OAuth2 tokens.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| data | body | Sign in credentials | Yes | [controllers.userLogin](#controllersuserlogin) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [controllers.userType](#controllersusertype) |
| 500 | Internal Server Error | [controllers.errorResult](#controllerserrorresult) |

### /users/modules/{id}

#### PUT
##### Summary

Update the module import link of a user

##### Description

Updates the module import link of the userId specified.

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path | userId | Yes | string |
| linktoadd | query | linkToAdd | Yes | string |
| token | header | Authorization token | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [controllers.userType](#controllersusertype) |
| 404 | Not Found | [controllers.errorResult](#controllerserrorresult) |

##### Security

| Security Schema | Scopes |
| --- | --- |
| ApiKeyAuth | |

### /users/signup

#### POST
##### Summary

User sign up

##### Description

Responds with userId

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| data | body | New user credentials | Yes | [controllers.userSignUp](#controllersusersignup) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [controllers.signUpResult](#controllerssignupresult) |
| 400 | Bad Request | [controllers.errorResult](#controllerserrorresult) |

### /users/update/{id}

#### PUT
##### Summary

Modify user particulars

##### Description

Change user particulars

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| id | path | userId | Yes | string |
| first_name | query | First name | No | string |
| last_name | query | Last name | No | string |
| email | query | Email | No | string |
| password | query | Password | No | string |
| token | header | Authorization token | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [controllers.userType](#controllersusertype) |
| 404 | Not Found | [controllers.errorResult](#controllerserrorresult) |

##### Security

| Security Schema | Scopes |
| --- | --- |
| ApiKeyAuth | |

### Models

#### controllers.errorResult

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| error | string |  | No |

#### controllers.signUpResult

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| InsertedID | string |  | No |

#### controllers.taskAddType

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| duration | string |  | No |
| first_name | string |  | Yes |
| hidden | boolean |  | No |
| last_name | string |  | Yes |
| moduleCode | string |  | No |
| taskName | string |  | No |
| user_id | string |  | No |

#### controllers.taskType

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| created_at | string |  | No |
| duration | string |  | No |
| first_name | string |  | Yes |
| hidden | boolean |  | No |
| id | string |  | No |
| last_name | string |  | Yes |
| moduleCode | string |  | No |
| taskName | string |  | No |
| updated_at | string |  | No |
| user_id | string |  | No |

#### controllers.userLogin

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| email | string |  | Yes |
| password | string |  | Yes |

#### controllers.userSignUp

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| email | string |  | No |
| first_name | string |  | Yes |
| last_name | string |  | Yes |
| password | string |  | No |

#### controllers.userType

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| Password | string |  | Yes |
| created_at | string |  | No |
| email | string |  | Yes |
| first_name | string |  | Yes |
| id | string |  | No |
| last_name | string |  | Yes |
| points | integer |  | No |
| refresh_token | string |  | No |
| timetable | string |  | No |
| token | string |  | No |
| updated_at | string |  | No |
| user_id | string |  | No |