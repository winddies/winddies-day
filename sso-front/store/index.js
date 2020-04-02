// import axios from 'axios'

// const url = require('url')

// export const actions = {
//   nuxtServerInit({ commit }, { req, redirect }) {
//     if (!req.headers.cookie) {
//       return
//     }
//     const cookies = req.headers.cookie.split(';')
//     const isLogin = cookies.some(item => (
//       item.trim().indexOf('SSO_ID') === 0
//     ))
//     if (isLogin) {
//       axios.get('http://0.0.0.0:8080/auth/verify-login', {
//         headers: { cookie: req.headers.cookie }
//       }).catch(err => console.error(err))
//       .then(res => {
//         const result = url.parse(req.url, true)
//         if (result.query && result.query.redirect) {
//           redirect(`${result.query.redirect}?token=${res.data.token}`)
//         }
//       })
//     }
//   }
// }
