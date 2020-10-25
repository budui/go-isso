import api from './api'
import $ from './dom'
import i18n from './i18n'

export default function () {
  var objs = {}

  $.each('a', function (el) {
    if (!el.href.match || !el.href.match(/#isso-thread$/)) {
      return
    }

    var tid =
      el.getAttribute('data-isso-id') ||
      el.href.match(/^(.+)#isso-thread$/)[1].replace(/^.*\/\/[^\/]+/, '')

    if (tid in objs) {
      objs[tid].push(el)
    } else {
      objs[tid] = [el]
    }
  })

  var urls = Object.keys(objs)

  api.count(urls).then(function (rv) {
    for (var key in objs) {
      if (objs.hasOwnProperty(key)) {
        var index = urls.indexOf(key)

        for (var i = 0; i < objs[key].length; i++) {
          objs[key][i].textContent = i18n.pluralize('num-comments', rv[index])
        }
      }
    }
  })
}
