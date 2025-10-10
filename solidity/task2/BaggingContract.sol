// SPDX-License-Identifier: MIT

pragma solidity ^0.8.4;

contract BaggingContract{
    //桶管理员
    address public owner;

    //记账本
    mapping(address => uint256) public donations;

    //放钱提醒事件
    event Donated(address indexed donor,uint256 amount);

    //取钱提醒时间
    event WithDrawn(address indexed owner,uint256 amount);

    //构造器
    constructor(){
        owner = msg.sender;
    }

    //修饰器
    modifier onlyOwner(){
        require(msg.sender == owner,"not the contract owner");
        _;
    }

    //放钱函数
    function donate() public payable {
        require(msg.value>0,"Must donate 1 at least");
        donations[msg.sender] += msg.value;
        emit Donated(msg.sender,msg.value);
    }


    //直接放钱函数
    receive() external payable {
        donate();
    }

    //取钱函数
    function withdraw() public onlyOwner{
        uint256 contructBalance = address(this).balance;
        require(contructBalance>0,"balance need to dayu 0");
        (bool success, ) = msg.sender.call{value: contructBalance}("");
        require(success,"transfer failed");
        emit WithDrawn(msg.sender,contructBalance);
    }

    //查账函数
    function getDonation(address donor) external view returns (uint256){
        return donations[donor];
    }


}