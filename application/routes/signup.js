var express = require('express');
var users = require('../public/js/users')
var cert = require('../public/js/cert')

const USER_COOKIE_KEY = 'USER'

var router = express.Router();

/* GET users listing. */
router.get('/', function(req, res, next) {
  res.render('signup', { title: 'CryptoFungi' });
});

// 일반 회원가입
router.post('/',async function(req, res, next) {
  const { username, name, password } = req.body;
  const job = "user"

  const exists = await users.fetchUser(username)
  if (exists) {
    res.status(400).send(`이미 존재하는 사용자입니다. : ${username}`);
  }

  const newUSer = {
    username,
    name,
    password,
    job,
  }
  // JSON 파일에 회원 정보 저장
  users.createUser(newUSer)

  // 사용자 인증서 생성 및 지갑저장
  cert.makeUsesrWallet(username, "org1")

  res.cookie(USER_COOKIE_KEY, JSON.stringify(newUSer))

  res.redirect("/")

});

// 먹이공장 회원가입
router.post('/feedfactory',async function(req, res, next) {
  const { username, name, password } = req.body;
  const job = "feedfactory"

  const exists = await users.fetchUser(username)
  if (exists) {
    res.status(400).send(`이미 존재하는 사용자입니다. : ${username}`);
    return
  }

  const newUSer = {
    username,
    name,
    password,
    job,
  }
  // JSON 파일에 회원 정보 저장
  users.createUser(newUSer)

  // 사용자 인증서 생성 및 지갑저장
  cert.makeUsesrWallet(username, "org2")

  res.cookie(USER_COOKIE_KEY, JSON.stringify(newUSer))

  res.redirect("/")
  
});

module.exports = router;
