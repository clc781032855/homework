import "openzeppelin-contracts/contracts/token/ERC20/IERC20.sol";

contract  MyContract is ERC20{
constructor(
    address initialOwner,
    uint256 initialSupply
)ERC20("Clinco", "Clc") Ownable(initialOwner){
    _mint(initialOwner, initialSupply*10**decimals());
}
}