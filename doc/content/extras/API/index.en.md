# API

The Isso API uses HTTP and JSON as primary communication protocol.

## JSON format

When querying the API you either get a regular HTTP error, an object or
list of objects representing the comment. Here's an example JSON
returned from Isso:

```json
{
    "id": 1,
    "parent": null,
    "text": "<p>Hello, World!</p>\n",
    "mode": 1,
    "hash": "4505c1eeda98",
    "author": null,
    "website": null,
    "created": 1387321261.572392,
    "modified": null,
    "likes": 3,
    "dislikes": 0
}
```

| Field    | Description                                                                                 |
| :------: | :-----------------------------------------------------------------------------------------: |
| id       | comment id (unique per website)                                                             |
| parent   | parent id reference, may be `null`                                                          |
| text     | required, comment written in Markdown                                                       |
| mode     | 1 – accepted       2 – in moderation queue    4 – deleted, but referenced.                  |
| hash     | user identication, used to generate identicons. PBKDF2 from email or IP address (fallback). |
| author   | author’s name, may be `null`                                                               |
| website  | author’s website, may be `null`                                                            |
| likes    | upvote count, defaults to 0                                                                 |
| dislikes | downvote count, defaults to 0                                                               |
| created  | time in seconds since UNIX time                                                             |
| modified | last modification since UNIX time, may be `null`                                            |

### comment mode description

| value | explanation                                                                                        |
| :---: | :------------------------------------------------------------------------------------------------: |
| 1     | accepted: The comment was accepted by the server and is published.                                 |
| 2     | in moderation queue: The comment was accepted by the server but awaits moderation.                 |
| 4     | deleted, but referenced: The comment was deleted on the server but is still referenced by replies. |

## fetch

----
  Queries the comments of a thread.

* **URL**

  `/`

  > /?uri=/thread/&limit=2&nested_limit=5

* **Method:**
  
  `GET`
  
* **URL Params**

| field        | type   | limit      | desc                                                                                    |
| :----------: | :----: | :--------: | :-------------------------------------------------------------------------------------: |
| uri          | string | `Required` | The URI of thread to get the comments from.                                             |
| parent       | number | `Optional` | Return only comments that are children of the comment with the provided ID.             |
| limit        | number | `Optional` | The maximum number of returned top-level comments. Omit for unlimited results.          |
| nested_limit | number | `Optional` | The maximum number of returned nested comments per commint. Omit for unlimited results. |
| after        | number | `Optional` | Includes only comments were added after the provided UNIX timestamp.                    |
| plain        | number | `Optional` | "0" for text which contains html tags,others for keep the text of comment plain.        |

* **Success Response:**

  * **Code:** 200

| field          | type     | desc                                                                                                                                                                                                  |
| :------------: | :------: | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------: |
| total_replies  | number   | The number of replies if the `limit` parameter was not set.<br>  If `after` is set to `X`, this is the number of comments that were created after `X`.<br>  So setting `after` may change this value! |
| replies        | Object[] | The list of comments. <br> Each comment also has the `total_replies`, `replies`, `id` <br> and `hidden_replies` properties to represent nested comments.                                              |
| id             | number   | Id of the comment `replies` is the list of replies of. <br> `null` for the list of toplevel comments.                                                                                                 |
| hidden_replies | number   | The number of comments that were ommited from the results <br> because of the `limit` request parameter.<br>  Usually, this will be `total_replies` - `limit`.                                        |

Example:

```output
$ curl 'https://comments.example.com/?uri=/thread/&limit=2&nested_limit=5'
{
  "total_replies": 14,
  "replies": [
    {
      "website": null,
      "author": null,
      "parent": null,
      "created": 1464818460.732863,
      "text": "&lt;p&gt;Hello, World!&lt;/p&gt;",
      "total_replies": 1,
      "hidden_replies": 0,
      "dislikes": 2,
      "modified": null,
      "mode": 1,
      "replies": [
        {
          "website": null,
          "author": null,
          "parent": 1,
          "created": 1464818460.769638,
          "text": "&lt;p&gt;Hi, now some Markdown: &lt;em&gt;Italic&lt;/em&gt;, &lt;strong&gt;bold&lt;/strong&gt;, &lt;code&gt;monospace&lt;/code&gt;.&lt;/p&gt;",
          "dislikes": 0,
          "modified": null,
          "mode": 1,
          "hash": "2af4e1a6c96a",
          "id": 2,
          "likes": 2
        }
      ],
      "hash": "1cb6cc0309a2",
      "id": 1,
      "likes": 2
    },
    {
      "website": null,
      "author": null,
      "parent": null,
      "created": 1464818460.80574,
      "text": "&lt;p&gt;Lorem ipsum dolor sit amet, consectetur adipisicing elit. Accusantium at commodi cum deserunt dolore, error fugiat harum incidunt, ipsa ipsum mollitia nam provident rerum sapiente suscipit tempora vitae? Est, qui?&lt;/p&gt;",
      "total_replies": 0,
      "hidden_replies": 0,
      "dislikes": 0,
      "modified": null,
      "mode": 1,
      "replies": [],
      "hash": "1cb6cc0309a2",
      "id": 3,
      "likes": 0
    },
    "id": null,
    "hidden_replies": 12
}
```

* **Error Response:**

| error                        | status code      | response                                     |
| :--------------------------: | :--------------: | :------------------------------------------: |
| can not find vaild comment   | `404` NotFound   | `{ "error" : "Not Found" }`                  |
| param `parent` invalid       | `400` BadRequest | `{ "error" : "param parent invalid" }`       |
| param `limit` invalid        | `400` BadRequest | `{ "error" : "param limit invalid" }`        |
| param `after` invalid        | `400` BadRequest | `{ "error" : "param after invalid" }`        |
| param `nested_limit` invalid | `400` BadRequest | `{ "error" : "param nested_limit invalid" }` |

## new

----
  The payload must be valid JSON. To prevent CSRF attacks, you must set
the `Content-Type` to `application/json` or omit the header
completely.

The server issues a cookie per new comment which acts as authentication
token to modify or delete your own comment. The token is
cryptographically signed and expires automatically after 900 seconds by
default.

The following keys can be used to POST a new comment, all other fields
are dropped or replaced with values from the server:

* **URL**

  `/new`

  > /new?uri=/thread/

* **Method:**
  
  `POST`
  
* **URL Params**

| field | type   | limit      | desc                                                                             |
| :---: | :----: | :--------: | :------------------------------------------------------------------------------: |
| uri   | string | `Required` | The URI of thread to get the comments from.                                      |
| plain | number | `Optional` | "0" for text which contains html tags,others for keep the text of comment plain. |

* **Payload Params**

| field   | type   | limit      | desc                                                                               |
| :-----: | :----: | :--------: | :--------------------------------------------------------------------------------: |
| text    | string | `Required` | The comment’s raw text.                                                           |
| author  | string | `Optional` | The comment’s author’s name.                                                     |
| email   | string | `Optional` | The comment’s author’s email address.                                            |
| website | string | `Optional` | The comment’s author’s website’s url.                                           |
| parent  | number | `Optional` | The parent comment’s id iff the new comment is a response to an existing comment. |

* **Success Response:**

  * **Code:** 202

  Return a typical comment object

Example:

```output
$ curl 'https://comments.example.com/new?uri=/thread/' -d '{"text": "Stop saying that! *isso*!", "author": "Max Rant", "email": "rant@example.com", "parent": 15}' -H 'Content-Type: application/json'

{
  "website": null,
  "author": "Max Rant",
  "parent": 15,
  "created": 1464940838.254393,
  "text": "&lt;p&gt;Stop saying that! &lt;em&gt;isso&lt;/em&gt;!&lt;/p&gt;",
  "dislikes": 0,
  "modified": null,
  "mode": 1,
  "hash": "e644f6ee43c0",
  "id": 23,
  "likes": 0
}
```

* **Error Response:**

| error                        | status code      | response                                     |
| :--------------------------: | :--------------: | :------------------------------------------: |
| param `parent` invalid       | `400` BadRequest | `{ "error" : "param parent invalid" }`       |
| param `limit` invalid        | `400` BadRequest | `{ "error" : "param limit invalid" }`        |
| param `after` invalid        | `400` BadRequest | `{ "error" : "param after invalid" }`        |
| param `nested_limit` invalid | `400` BadRequest | `{ "error" : "param nested_limit invalid" }` |

## Edit comment

When your authentication token is not yet expired, you can issue a PUT
request to update `text`, `author` and `website`. After an update, you get an
updated authentication token and the comment as JSON:

``` bash
curl -X PUT http://isso/id/1 -d "..." -H "Content-Type: application/json"
```

## Delete comment

You can delete your own comments when your authentication token (=cookie) is not yet expired:

```bash
$ curl -X DELETE http://isso/id/1 -H "Content-Type: application/json"
null
```

Returns either `null` or a comment with an empty text value when the comment is still referenced by other comments.

## Up and downvote comments

...

## Get comment count

Counts all publicly visible comments for thread `uri`:

```text
GET /count?uri=%2Fhello-world%2F
2
```

| Field | Description                          |
| :---: | :----------------------------------: |
| uri   | URI to count comments for, required. |

returns an integer

## Get Atom feed

Get an Atom feed of comments for thread `uri`:

```text
GET /feed?uri=%2Fhello-world%2F
```

| Field | Description                          |
| :---: | :----------------------------------: |
| uri   | URI to count comments for, required. |

Returns an XML document as the Atom feed.