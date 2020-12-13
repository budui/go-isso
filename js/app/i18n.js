import config from './config'
import bg from './i18n/bg'
import cs from './i18n/cs'
import da from './i18n/da'
import de from './i18n/de'
import en from './i18n/en'
import fa from './i18n/fa'
import fi from './i18n/fi'
import fr from './i18n/fr'
import hr from './i18n/hr'
import hu from './i18n/hu'
import ru from './i18n/ru'
import it from './i18n/it'
import eo from './i18n/eo'
import oc from './i18n/oc'
import pl from './i18n/pl'
import pt_BR from './i18n/pt_BR'
import sk from './i18n/sk'
import sv from './i18n/sv'
import nl from './i18n/nl'
import el_GR from './i18n/el_GR'
import es from './i18n/es'
import vi from './i18n/vi'
import zh_CN from './i18n/zh_CN'
import zh_TW from './i18n/zh_TW'

var pluralforms = function (lang) {
  switch (lang) {
    case 'bg':
    case 'cs':
    case 'da':
    case 'de':
    case 'el':
    case 'el_GR':
    case 'en':
    case 'es':
    case 'eo':
    case 'fa':
    case 'fi':
    case 'hr':
    case 'hu':
    case 'it':
    case 'pt_BR':
    case 'sv':
    case 'nl':
    case 'vi':
    case 'zh':
    case 'zh_CN':
    case 'zh_TW':
      return function (msgs, n) {
        return msgs[n === 1 ? 0 : 1]
      }
    case 'fr':
      return function (msgs, n) {
        return msgs[n > 1 ? 1 : 0]
      }
    case 'ru':
      return function (msgs, n) {
        if (n % 10 === 1 && n % 100 !== 11) {
          return msgs[0]
        } else if (
          n % 10 >= 2 &&
          n % 10 <= 4 &&
          (n % 100 < 10 || n % 100 >= 20)
        ) {
          return msgs[1]
        } else {
          return typeof msgs[2] !== 'undefined' ? msgs[2] : msgs[1]
        }
      }
    case 'oc':
      return function (msgs, n) {
        return msgs[n > 1 ? 1 : 0]
      }
    case 'pl':
      return function (msgs, n) {
        if (n === 1) {
          return msgs[0]
        } else if (
          n % 10 >= 2 &&
          n % 10 <= 4 &&
          (n % 100 < 10 || n % 100 >= 20)
        ) {
          return msgs[1]
        } else {
          return typeof msgs[2] !== 'undefined' ? msgs[2] : msgs[1]
        }
      }
    case 'sk':
      return function (msgs, n) {
        if (n === 1) {
          return msgs[0]
        } else if (n === 2 || n === 3 || n === 4) {
          return msgs[1]
        } else {
          return typeof msgs[2] !== 'undefined' ? msgs[2] : msgs[1]
        }
      }
    default:
      return null
  }
}

// useragent's prefered language (or manually overridden)
var lang = config.lang

// fall back to English
if (!pluralforms(lang)) {
  lang = 'en'
}

var catalogue = {
  bg: bg,
  cs: cs,
  da: da,
  de: de,
  el: el_GR,
  el_GR: el_GR,
  en: en,
  eo: eo,
  es: es,
  fa: fa,
  fi: fi,
  fr: fr,
  it: it,
  hr: hr,
  hu: hu,
  oc: oc,
  pl: pl,
  pt: pt_BR,
  pt_BR: pt_BR,
  ru: ru,
  sk: sk,
  sv: sv,
  nl: nl,
  vi: vi,
  zh: zh_CN,
  zh_CN: zh_CN,
  zh_TW: zh_TW
}

var plural = pluralforms(lang)

var translate = function (msgid) {
  return (
    config[msgid + '-text-' + lang] ||
    catalogue[lang][msgid] ||
    en[msgid] ||
    '???'
  )
}

var pluralize = function (msgid, n) {
  var msg

  msg = translate(msgid)
  if (msg.indexOf('\n') > -1) {
    msg = plural(msg.split('\n'), +n)
  }

  return msg ? msg.replace('{{ n }}', +n) : msg
}

export default {
  lang: lang,
  translate: translate,
  pluralize: pluralize
}
