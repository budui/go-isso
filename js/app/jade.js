import utils from './utils'
import tt_comment from './text/comment'
import tt_comment_loader from './text/comment_loader'
import tt_postbox from './text/postbox'

var globals = {}
var templates = {
  comment: tt_comment,
  comment_loader: tt_comment_loader,
  postbox: tt_postbox
}

var set = function (name, value) {
  globals[name] = value
}

set('bool', function (arg) {
  return !!arg
})
set('humanize', function (date) {
  if (typeof date !== 'object') {
    date = new Date(parseInt(date, 10) * 1000)
  }

  return date.toString()
})
set('datetime', function (date) {
  if (typeof date !== 'object') {
    date = new Date(parseInt(date, 10) * 1000)
  }

  return (
    [
      date.getUTCFullYear(),
      utils.pad(date.getUTCMonth(), 2),
      utils.pad(date.getUTCDay(), 2)
    ].join('-') +
    'T' +
    [
      utils.pad(date.getUTCHours(), 2),
      utils.pad(date.getUTCMinutes(), 2),
      utils.pad(date.getUTCSeconds(), 2)
    ].join(':') +
    'Z'
  )
})

export default {
  set: set,
  render: function (name, locals) {
    var rv,
      t = templates[name]
    if (!t) {
      throw new Error("Template not found: '" + name + "'")
    }

    locals = locals || {}

    var keys = []
    for (var key in locals) {
      if (locals.hasOwnProperty(key) && !globals.hasOwnProperty(key)) {
        keys.push(key)
        globals[key] = locals[key]
      }
    }

    rv = templates[name](globals)

    for (var i = 0; i < keys.length; i++) {
      delete globals[keys[i]]
    }

    return rv
  }
}
