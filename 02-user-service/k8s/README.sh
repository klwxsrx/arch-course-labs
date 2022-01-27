# Install mysql from helm before:
helm install --set "image.tag=8.0" --set "auth.database=user" --set "auth.username=user" --set "auth.password=test1234" userservice-db bitnami/mysql