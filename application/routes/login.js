var express = require('express');
var users = require('../public/js/users')

var router = express.Router();
const USER_COOKIE_KEY = 'USER'

/* GET users listing. */
router.get('/', function(req, res, next) {
  res.render('login', { title: 'CryptoFungi' });
});

router.post('/', async function(req, res, next) {
    const { username, password } = req.body;
    const user = await users.fetchUser(username)

    if(!user) {
        res.status(400).send(`가입되지 않은 사용자입니다. : ${username}`);
        return;
    }

    if (password !== user.password) {
        res.status(400).send(`비밀번호가 틀렸습니다.`);
        return;
    }
    res.cookie(USER_COOKIE_KEY, JSON.stringify(user))

    res.redirect("/")
});

module.exports = router;