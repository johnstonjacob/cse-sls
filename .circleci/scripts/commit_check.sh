#! /bin/bash
set -e

# latest commit
LATEST_COMMIT=$(git rev-parse HEAD)

# latest commit where path/to/folder1 was changed
FOLDER1_COMMIT=$(git log -1 --format=format:%H --full-diff ~/proj/workflow-cost-estimator)

# latest commit where path/to/folder2 was changed
#FOLDER2_COMMIT=$(git log -1 --format=format:%H --full-diff path/to/folder2)

if [ $FOLDER1_COMMIT = $LATEST_COMMIT ];
    then
        echo "files in folder1 has changed"
        ~/proj/.circleci/scripts/workflow-cost-estimator.sh
#elif [ $FOLDER2_COMMIT = $LATEST_COMMIT ];
#    then
#        echo "files in folder2 has changed"
#        .circleci/do_something_else.sh
else
     echo "no folders of relevance has changed"
     exit 0;
 fi
