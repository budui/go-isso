# roadmap for go-isso

## Port

First, I need to implement all the functions of `isso` with golang.


| Status | Milestone  | Goals |
| :----: | :------------------------ | :---: |
|   ❌   | **[Frontend](#Frontend)** | 0 / ∞ |
|   🚀   | **[Database](#Database)** | 14 / 22 |
|   🚀   | **[API](#API)** | 3 / 17 |
|   ❌   | **[Notifications](#Notifications)** | 0 / 2 |
|   ❌   | **[Guard](#API)** | 0 / 2 |

### Frontend

I wont touch the front part before I finish all backend stuff.

### Database

the database have 3 table: **threads**, **preferences** and **comments**.

#### threads

| name |  status |
| ---- | ---- |
|  NewThread    | ✔  |
|  GetThreadWithID    | ✔ |
|  GetThreadWithUri    | ✔ |
|  Contain    | ✔ |

#### preferences

| name           | status |
| -------------- | ------ |
| initPreference | ✔      |
| GetPreference  | ✔      |
| SetPreference  | ✔      |

#### comments

| name          | status |
| ------------- | ------ |
| _remove_stale | ❌      |
| init          | ✔      |
| add           | ✔      |
| activate      | ❌      |
| unsubscribe   | ❌      |
| fetchall      | ❌      |
| fetch         | ✔      |
| delete        | ✔      |
| update        | ✔      |
| get           | ✔      |
| count_modes   | ❌      |
| vote          | ❌      |
| reply_count   | ✔      |
| count         | ❌      |
| purge         | ❌      |

### API

| name        | route                                                        | status |
| ----------- | ------------------------------------------------------------ | ------ |
| fetch       | (`GET`, /)                                                   | ✔      |
| new         | (`POST`, /new)                                               | ✔      |
| count       | (`GET`, /count)                                              | ❌      |
| counts      | (`POST`, /count)                                             | ❌      |
| feed        | (`GET`, /feed)                                               | ❌      |
| view        | (`GET`, /id/\<int:id\>)                                      | ✔      |
| edit        | (`PUT`, /id/\<int:id\>)                                      | 🚀      |
| delete      | (`DELETE`, /id/\<int:id\>)                                   | ❌      |
| unsubscribe | (`GET`, /id/\<int:id\>/unsubscribe/\<string:email\>/\<string:key\>) | ❌      |
| moderate    | (`GET`,  /id/\<int:id\>/\<any(edit,activate,delete):action\>/\<string:key\>) | ❌      |
| moderate    | (`POST`, /id/\<int:id\>/\<any(edit,activate,delete):action\>/\<string:key\>) | ❌      |
| like        | (`POST`, /id/\<int:id\>/like)                                | ❌      |
| dislike     | (`POST`, /id/\<int:id\>/dislike)                             | ❌      |
| demo        | (`GET`, /demo)                                               | ❌      |
| preview     | (`POST`, /preview)                                           | ❌      |
| login       | (`POST`, /login)                                             | ❌      |
| admin       | (`GET`, /admin)                                              | ❌      |

### Guard

| name       | status |
| ---------- | ------ |
| race limit | ❌      |
| spam       | ❌      |

### Notifications

| name  | status |
| ----- | ------ |
| email | ❌      |
| log   | ❌      |

## Beyond isso

1. Telegram support
2. optional 3-party service support