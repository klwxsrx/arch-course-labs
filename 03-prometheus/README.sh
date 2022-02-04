# Install prometheus
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
helm install prom prometheus-community/kube-prometheus-stack -f ./kube-prometheus-stack.values.yaml --atomic --namespace monitoring --create-namespace

# Install ingress controller with enabled metrics
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo update
helm install nginx ingress-nginx/ingress-nginx -f ingress-nginx.values.yaml --atomic --namespace ingress-nginx --create-namespace

# Install user service
helm install user-service ../02-user-service/helm/ --namespace user-service --create-namespace

# Install grafana dashboard for the service
kubectl apply -f ./grafana.yaml