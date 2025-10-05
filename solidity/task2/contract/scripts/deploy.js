const { ethers } = require("hardhat");

async function main() {
    const [deployer] = await ethers.getSigners();

    console.log("使用账户部署合约:", deployer.address);
    console.log("账户余额:", (await deployer.provider.getBalance(deployer.address)).toString());

    // 获取合约工厂
    const MyToken = await ethers.getContractFactory("MyToken");

    // 部署合约 - 初始供应1000个代币
    const initialSupply = ethers.parseEther("1000");

    console.log("正在部署 MyToken 合约...");
    const myToken = await MyToken.deploy(deployer.address, initialSupply);

    await myToken.waitForDeployment();

    const contractAddress = await myToken.getAddress();
    console.log("✅ MyToken合约部署成功!");
    console.log("📝 合约地址:", contractAddress);
    console.log("👤 合约所有者:", deployer.address);
    console.log("💰 初始供应量:", initialSupply.toString());

    const totalSupply = await myToken.totalSupply();
    console.log("总供应量:", totalSupply.toString());
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error("❌ 部署失败:", error);
        process.exit(1);
    });