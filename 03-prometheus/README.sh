# Install prometheus
kubectl create namespace monitoring && kubectl config set-context --current --namespace=monitoring
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo add stable https://charts.helm.sh/stable
helm repo update
helm install prom prometheus-community/kube-prometheus-stack -f ./kube-prometheus-stack.values.yaml --atomic

# Install user service
kubectl create namespace user-service && kubectl config set-context --current --namespace=user-service
helm install user-service ../02-user-service/helm/