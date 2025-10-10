pragma solidity ^0.8.20;

// 导入OpenZeppelin的ERC20实现
import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/access/Ownable.sol";



contract MyNFT is ERC721, Ownable {

    uint256 private _tokenIdCounter;

    constructor(
    ) ERC721("MyNFT", "MNFT") Ownable(msg.sender) {
    }

    function mintNFT(address recipient, string memory tokenURI) public onlyOwner returns (uint256){
        _tokenIdCounter++;
        uint256 tokenId = _tokenIdCounter;
        _safeMint(recipient, tokenId);
        _setTokenURI(tokenId, tokenURI);
        return tokenId;
    }

    function tokenURI(uint256 tokenId) public view override returns (string memory) {
        return super.tokenURI(tokenId);
    }
}
