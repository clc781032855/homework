// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

// 导入OpenZeppelin的ERC20实现
import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

/**
 * @title 我的ERC20代币
 * @dev 这是一个简单的ERC20代币实现，包含增发功能
 */
contract MyToken is ERC20, Ownable {
    /**
     * @dev 构造函数，初始化代币
     * @param initialOwner 合约所有者地址
     * @param initialSupply 初始供应量
     */
    constructor(
        address initialOwner,
        uint256 initialSupply
    ) ERC20("MyToken", "MTK") Ownable(initialOwner) {
        _mint(initialOwner, initialSupply * 10 ** decimals());
    }

    /**
     * @dev 增发代币（仅所有者可调用）
     * @param to 接收地址
     * @param amount 增发金额
     */
    function mint(address to, uint256 amount) public onlyOwner {
        _mint(to, amount);
    }
}