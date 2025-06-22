#!/bin/bash
set -e

echo "📦 Installing Kind (Kubernetes in Docker)..."

# Detect OS
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
    x86_64) ARCH="amd64" ;;
    arm64|aarch64) ARCH="arm64" ;;
    *) echo "❌ Unsupported architecture: $ARCH"; exit 1 ;;
esac

# Download and install Kind
KIND_VERSION="v0.20.0"
KIND_URL="https://kind.sigs.k8s.io/dl/${KIND_VERSION}/kind-${OS}-${ARCH}"

echo "🔽 Downloading Kind ${KIND_VERSION} for ${OS}-${ARCH}..."
curl -Lo ./kind "${KIND_URL}"
chmod +x ./kind

# Install to system path
echo "Installing Kind to /usr/local/bin (requires sudo)..."
if [[ "$OS" == "linux" ]] || [[ "$OS" == "darwin" ]]; then
    sudo mv ./kind /usr/local/bin/kind
    sudo chmod +x /usr/local/bin/kind
else
    echo "❌ Unsupported OS: $OS"
    exit 1
fi

# Verify installation
if command -v kind &> /dev/null; then
    echo "✅ Kind installed successfully!"
    kind version
else
    echo "❌ Kind installation failed"
    exit 1
fi

# Check kubectl
if ! command -v kubectl &> /dev/null; then
    echo "⚠️  kubectl not found. Installing..."
    
    # Install kubectl
    KUBECTL_VERSION=$(curl -L -s https://dl.k8s.io/release/stable.txt)
    KUBECTL_URL="https://dl.k8s.io/release/${KUBECTL_VERSION}/bin/${OS}/${ARCH}/kubectl"
    
    curl -LO "${KUBECTL_URL}"
    chmod +x kubectl
    sudo mv kubectl /usr/local/bin/
    sudo chmod +x /usr/local/bin/kubectl
    
    echo "✅ kubectl installed successfully!"
    kubectl version --client
fi

echo ""
echo "🎉 Kind and kubectl are ready!"
echo "💡 You can now run Kafka setup with: ./setup_strimzi.sh setup"