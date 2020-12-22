from example_chaincode import OperateChaincode
from hfc.fabric import Client

cli = Client(net_profile="example_server_network.json")
org1_admin = cli.get_user('org1.example.com', 'Admin')

# definition operator and chaincode
OC = OperateChaincode(org1_admin, 'mychannel',['peer0.org1.example.com'], 'mycc')

#transactionID := args[0]
args = ("0000000000100000")
# response the query
resp = OC.groupbuyBlocked_failed([*args])
print(resp)



