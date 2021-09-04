const hash = window.location.hash.substring(1).split('&')
var datas = {}
for (let elem of hash) {
  elem = elem.split('=')
  datas[elem[0]] = elem[1]
}
fetch('/mytoken', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify(datas)
}).then(() => {window.location.replace("/")})
