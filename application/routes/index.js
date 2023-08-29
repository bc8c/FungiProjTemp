var express = require('express');
var router = express.Router();
var cc = require('../public/js/cc')

const USER_COOKIE_KEY = 'USER'

/* GET home page. */
router.get('/', async function(req, res, next) {
  const userCookie = req.cookies[USER_COOKIE_KEY];
  console.log(userCookie)

  if (!userCookie) {  // 로그인이 안된경우
    res.render('users', {title:"CryptoFungi"})
  } else {            // 로그인이 되어있는 경우
    
    const userData = JSON.parse(userCookie)
    const id = userData.username    
    console.log(id)
    // 1. 일반사용자
    if (userData.job == "user"){
      var result = await cc.cc_call(id, "GetFungiByOwner", "")
      res.render('index', { title: 'CryptoFungi', result:result });
      console.log(result)
    }
    else { // 2. 먹이공장직원
      res.render('index', { title: 'CryptoFungi', result:"feed" });
    }
  }  
});

router.get('/logout', function(req, res, next) {
  res.clearCookie(USER_COOKIE_KEY);
  res.redirect('/')
});

module.exports = router;
