sed -e "/$2:/,/.*/s|image: .*|image: ${CI_REGISTRY_IMAGE}:$2-${CI_COMMIT_REF_NAME}-${CI_COMMIT_SHORT_SHA}|" $1 > $1.tmp
mv $1.tmp $1