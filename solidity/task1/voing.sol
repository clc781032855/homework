pragma solidity ^0.8.0;

contract Voting{
    mapping(string => unit256) private votes;

    string[] private candidates;

    // 事件：当用户投票时触发
    event Voted(address indexed voter, string indexed candidate);

    // 事件：当所有候选人的得票数重置时触发
    event VotesReset();

    // 构造函数：初始化候选人列表
    constructor(string[] memory _candidates) {
        candidates = _candidates;
        for (uint i = 0; i < _candidates.length; i++) {
            votes[_candidates[i]] = 0;
        }
    }

    // 投票函数：允许用户投票给某个候选人
    function vote(string memory candidate) public{
        require(isValidCandidate(candidate),"Invalid candidate");
        votes[candidate]++;
        emit Voted(msg.sender, candidate);
    }

    // 获取候选人得票数
    function getVoteCount(string memory candidate) public view returns (uint256) {
        require(isValidCandidate(candidate), "Invalid candidate");
        return votes[candidate];
    }

    // 重置所有候选人的得票数
    function resetVotes() public {
        votes["A"] = 0;
        votes["B"] = 0;
        emit VotesReset();
    }

    // 检查候选人是否有效
    function isValidCandidate(string memory candidate) private view returns (bool) {
        for (uint i = 0; i < candidates.length; i++) {
            if (keccak256(bytes(candidates[i])) == keccak256(bytes(candidate))) {
                return true;
            }
        }
        return false;
    }

    // 获取所有候选人列表
    function getCandidates() public view returns (string[] memory) {
        return candidates;
    }
}