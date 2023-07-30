rm ./out -r
mkdir out
python src/sssg.py -i ./example -o ./out
pushd out
python -m http.server 8000 --bind 127.0.0.1
popd

