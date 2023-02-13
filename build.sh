COMMIT_ID=`git log | head -1 | sed s/'commit '//`
SUB_COMMIT_ID="${COMMIT_ID:0:7}"
STR_DATE=`git log | head -3 | tail -1`
COMMIT_DATE=`git log | head -3 | tail -1 | sed s/'Date:   '//`
if [[ $STR_DATE != *"Date:"* ]]; then
    COMMIT_DATE=`git log | head -4 | tail -1 | sed s/'Date:   '//`
fi
COMMIT_DATE_FORMAT=`date -r $(LANG=C date -juf '%a %b %d %T %Y %z' "$COMMIT_DATE" +%s) +'%y.%m%d.%H%M'`
IMAGE="go-redis-streaming:${COMMIT_DATE_FORMAT}_${SUB_COMMIT_ID}"
echo "$IMAGE"

docker build --platform linux/amd64 -t $IMAGE .