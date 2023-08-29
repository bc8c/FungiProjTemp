# 크립토풍기 프로젝트 실행 방법

## 1. BlockChain Network 실행

./startNetwork.sh

## 2. 체인코드 배포
### 먹이공장(feedfactory)
./deployFeedCC.sh 

### 버섯공장(fungusfactory)
./deployFungiCC.sh

## 3. 체인코드 테스트 ( Initalize를 위해 반드시 수행필요 )
### 먹이공장(feedfactory)
./testFeedCC.sh

### 버섯공장(fungusfactory)
./testFungiCC.sh

## 4. Application 실행
./startApplication.sh

## 5. Web Client 접속
웹브라우저 : localhost:3000

## END : 네트워크 종료
./downNetwork.sh