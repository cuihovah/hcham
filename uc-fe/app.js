const express = require('express')
const path = require('path')
const app = express()
const cookieParser = require('cookie-parser')
const ejs = require('ejs')
const config = require('config')
const url = require('url')
const JSONRPC = require('node-go-jsonrpc')
const bodyParser = require('body-parser')

const port = config.get('port')
const sso = new JSONRPC(config.get('sso.rpc_ip'), config.get('sso.rpc_port'))
const uc = new JSONRPC(config.get('uc.rpc_ip'), config.get('uc.rpc_port'))

app.engine('.html', ejs.__express)
app.set('view engine', 'html');
app.set('views', path.join(__dirname, 'views'));

app.use(cookieParser())
app.use(bodyParser.json());
app.use(bodyParser.urlencoded({extended: true}));

app.get('/register', async function (req, res) {
    res.render('register')
})

app.get('/user_api/rsa', async function (req, res) {
    const rsa = await uc.call('UCServer.RPCRSA', [])
    res.setHeader('Content-Type', 'text/plain')
    res.status(200).send(rsa.result)
})

app.post('/user_api/users', async function (req, res) {
    try {
        const ret = await uc.call('UCServer.RPCRegister', [req.body])
        res.status(201).end()
    } catch (err) {
        
    }
})

app.use(async function(req, res, next){
    const name = config.get('session.name')
    if (req.cookies[name] == undefined && req.query.token == undefined) {
        const host = config.get('sso.url')
        const self = config.get('self.url')
        const r = `${host}?redirect=${encodeURIComponent(self)}`
        return res.redirect(r)
    }
    if (req.query.token) {
        const ret = await sso.call('SSOServ.DecodeToken', [req.query.token])
        req.user = ret.result
        const cookie = JSON.stringify(ret.result)
        res.cookie(name, encodeURIComponent(cookie), {
            path: '/',
            domain: config.get('session.domain')
        })
    }
    if (req.cookies[name]) {
        let ustr = decodeURIComponent(req.cookies[name])
        const user = JSON.parse(ustr)
        req.user = user
    }
    return next()
})

app.get('/', async function (req, res) {
    const info = await uc.call('UCServer.RPCSearch', [req.user.id])
    res.render('update-user', info.result)
})

app.put('/user_api/users', async function (req, res) {
    const body = req.body
    body.id = req.user.id
    try {
        const ret = await uc.call('UCServer.RPCUpdate', [body])
        if (ret.result == true) {
            return res.status(200)
        } else {
            return res.status(400)
        }
    } catch(err) {
        return res.status(500)
    }
})

app.listen(process.env.NODE_PORT || port)