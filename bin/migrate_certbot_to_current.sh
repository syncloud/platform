#!/bin/bash -e

OLD_CERTBOT=/var/snap/platform/common/certbot
NEW_CERTBOT=/var/snap/platform/current/certbot

# fix current migration
if [[ -d $OLD_CERTBOT && -d $NEW_CERTBOT ]]; then
  certs=$(ls -la $NEW_CERTBOT/live | wc -l)
  if [[ $certs -gt 5 ]]; then
    echo "multiple certs detected, redoing the migration"
    rm -rf $NEW_CERTBOT
  fi
fi

# migrate common to current
if [[ ! -d $NEW_CERTBOT ]]; then
  cp -r $OLD_CERTBOT $NEW_CERTBOT
  RENEWAL_DIR=$NEW_CERTBOT/renewal
  if ls -la $RENEWAL_DIR/*.conf; then
      sed -i 's#'$OLD_CERTBOT'#'$NEW_CERTBOT'#g' $RENEWAL_DIR/*.conf
  fi
fi