./build.sh
git add .
git commit --amend -m ""
eb deploy
eb status
