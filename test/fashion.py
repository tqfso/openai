import os
import torch
from datetime import datetime
from torch import nn
from torch.utils.data import DataLoader
from torchvision import datasets
from torchvision.transforms import ToTensor

os.environ['HF_ENDPOINT'] = 'https://hf-mirror.com'

current_file_dir = os.path.dirname(os.path.abspath(__file__))
models_path = f'{current_file_dir}/models'
dataset_path = f'{current_file_dir}/data'

# 定义模型
class NeuralNetwork(nn.Module):
    def __init__(self):
        super().__init__()
        self.flatten = nn.Flatten()
        self.linear_relu_stack = nn.Sequential(
            nn.Linear(28*28, 512),
            nn.ReLU(),
            nn.Linear(512, 512),
            nn.ReLU(),
            nn.Linear(512, 10)
        )
    
    def forward(self, x):
        x = self.flatten(x)
        logits = self.linear_relu_stack(x)
        return logits


def _dataloader():

    global train_data, test_data
    global train_dataloader, test_dataloader

    # download training data from open datasets
    train_data = datasets.FashionMNIST(
        root=dataset_path, 
        train=True, 
        download=True, 
        transform=ToTensor())
    
    # Download test data from open datasets.
    test_data = datasets.FashionMNIST(
        root=dataset_path,
        train=False,
        download=True,
        transform=ToTensor(),
    )

    train_dataloader = DataLoader(
        dataset=train_data, 
        batch_size=64,        
        num_workers=0, 
        prefetch_factor=None,
        persistent_workers=False,
        pin_memory=True if device != 'cpu' else False,
        pin_memory_device="" if device == 'cpu' else device,        
        shuffle=True)
    
    test_dataloader = DataLoader(
        dataset=test_data, 
        batch_size=64, 
        num_workers=0, 
        prefetch_factor=None,
        persistent_workers=False,
        pin_memory=True if device != 'cpu' else False,
        pin_memory_device="" if device == 'cpu' else device,        
        shuffle=True,)

    for X, y in train_dataloader:
        print(f"Shape of X [N, C, H, W]: {X.shape}") # torch.Size([64, 1, 28, 28])
        print(f"Shape of y: {y.shape} {y.dtype}") # torch.Size([64]) torch.int64
        break

def _train(dataloader:DataLoader, model:nn.Module, loss_fn:nn.CrossEntropyLoss, optimizer:torch.optim.Optimizer):

    size = len(dataloader.dataset)
    model.train() # 训练模式
    for batch, (X, y) in enumerate(dataloader):
        X, y = X.to(device), y.to(device)

        # Compute prediction error
        pred = model(X)
        loss = loss_fn(pred, y)

        # Backpropagation
        loss.backward()
        optimizer.step()
        optimizer.zero_grad()

        if batch % 100 == 0:
            loss, current = loss.item(), (batch + 1) * len(X)
            print(f"loss: {loss:>7f}  [{current:>5d}/{size:>5d}]")

def _test(dataloader:DataLoader, model:nn.Module, loss_fn:nn.CrossEntropyLoss, ):
    size = len(dataloader.dataset)
    num_batches = len(dataloader)
    model.eval()
    test_loss, correct = 0, 0
    with torch.no_grad():
        for X, y in dataloader:
            X, y = X.to(device), y.to(device)
            pred = model(X)
            test_loss += loss_fn(pred, y).item()
            correct += (pred.argmax(1) == y).float().sum().item()
    test_loss /= num_batches
    correct /= size
    print(f"Test Error: \n Accuracy: {(100*correct):>0.1f}%, Avg loss: {test_loss:>8f} \n")

def train():

    model = NeuralNetwork().to(device=device)
    print(model)

    loss_fn = nn.CrossEntropyLoss() # 损失函数
    optimizer = torch.optim.SGD(model.parameters(), lr=1e-3) # 优化器
    
    _dataloader()

    epochs = 50
    for t in range(epochs):
        print(f"Epoch {t+1}\n-------------------------------")
        _train(train_dataloader, model, loss_fn, optimizer)
        _test(test_dataloader, model, loss_fn)
    
    print("Done!")

    # timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")  # 例如: 20251225_143022
    # file_path = os.path.abspath(f'{models_path}/fashion.pth')
    # model_dir = os.path.dirname(file_path)
    # os.makedirs(model_dir, exist_ok=True)
    # torch.save(model.state_dict(), file_path)
    # print(f"Saved PyTorch Model State to {file_path}")

def predict():

    _dataloader()
    
    file_path = os.path.abspath(f'{models_path}/fashion.pth')
    model = NeuralNetwork().to(device=device)
    model.load_state_dict(torch.load(file_path, weights_only=True))

    classes = [
        "T-shirt/top",
        "Trouser",
        "Pullover",
        "Dress",
        "Coat",
        "Sandal",
        "Shirt",
        "Sneaker",
        "Bag",
        "Ankle boot",
    ]

    model.eval()
    x, y = test_data[0][0], test_data[0][1]
    with torch.no_grad():
        x = x.unsqueeze(0)  # 形状将是 [1, C, H, W]
        x = x.to(device)
        pred = model(x)
        probes = nn.Softmax(dim=1)(pred)
        index = probes.argmax(1).item()

        predicted, actual = classes[index], classes[y]
        print(f'Predicted: "{predicted}", Actual: "{actual}"')

if __name__ == '__main__':
    device = 'cuda' if torch.cuda.is_available() else "cpu"
    print(f"Using {device} device")
    train()
    predict()
    print("Exit")
