var express = require('express');
var cc = require('../public/js/cc')

var router = express.Router();
const USER_COOKIE_KEY = 'USER'

/* GET users listing. */
router.get('/', function(req, res, next) {
  res.render('createFungus', { title: 'CryptoFungi' });
});


router.post('/', async function(req, res, next) {
    const name = req.body.name;
    const userCookie = req.cookies[USER_COOKIE_KEY];
    const userData = JSON.parse(userCookie);

    const id = userData.username    
    console.log(id)
    
    // 버섯 생성 부분
    var result = await cc.cc_call(id, "CreateRandomFungus", name)
    console.log(result);
    res.redirect("/");
});

router.post('/feed', async function(req, res, next) {
  const userCookie = req.cookies[USER_COOKIE_KEY];
  const userData = JSON.parse(userCookie);

  const id = userData.username    
  const fungusid = req.body.fungusid;
  const feedid = req.body.feedid;

  var args = [fungusid, feedid]
  
  // 버섯 생성 부분
  var result = await cc.cc_call(id, "Feed", args)
  console.log(result);
  res.redirect("/");
});

module.exports = router;