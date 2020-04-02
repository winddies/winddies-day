const Koa = require('koa')
const consola = require('consola')
const { Nuxt, Builder } = require('nuxt')
const request = require('request');
const cors = require('koa2-cors')
const url = require('url')

const app = new Koa()

// Import and Set Nuxt.js options
let config = require('../nuxt.config.js')
config.dev = !(app.env === 'production')

app.use(cors())

async function start() {
  // Instantiate nuxt.js
  const nuxt = new Nuxt(config)

  const {
    host = process.env.HOST || '127.0.0.1',
    port = process.env.PORT || 3000
  } = nuxt.options.server

  // Build in development
  if (config.dev) {
    const builder = new Builder(nuxt)
    await builder.build()
  } else {
    await nuxt.ready()
  }

  function getTokenRes(ctx) {
    // TODO: 更优雅的调用方式，比如通过 rpc
    return new Promise((resolve, reject) => {
      const options = {
        url: 'http://0.0.0.0:8080/auth/verify-login',
        headers: { cookie: ctx.req.headers.cookie }
      }
      request(options, function(err, response, body) {
        if (!err && response.statusCode == 200) {
          const token = JSON.parse(body)["token"]
          resolve(token)
        }
      })
    })
  }

  // 对于不带 cookie 或者 path 不是根的不用去请求 sso 判断是否登陆
  app.use(async (ctx, next) => {
    if (ctx.request.path !== '/' || !ctx.request.headers.cookie) {
      next()
      return
    }

    const cookies = ctx.req.headers.cookie.split(';')

    const isLogin = cookies.some(item => (
      item.trim().indexOf('SSO_ID') === 0
    ))

    const result = url.parse(ctx.req.url, true)

    if (isLogin && result.query && result.query.redirect) {
      const token = await getTokenRes(ctx)
      ctx.redirect(`${result.query.redirect}?token=${token}`)
      return
    }
    next()
  })

  app.use(ctx => {
    ctx.status = 200
    ctx.respond = false // Bypass Koa's built-in response handling
    ctx.req.ctx = ctx // This might be useful later on, e.g. in nuxtServerInit or with nuxt-stash
    nuxt.render(ctx.req, ctx.res)
  })

  app.listen(port, host)
  consola.ready({
    message: `Server listening on http://${host}:${port}`,
    badge: true
  })
}

start()
