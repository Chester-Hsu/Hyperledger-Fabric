import asyncio
from hfc.fabric import Client
import os
import json

loop = asyncio.get_event_loop()
dir_path = os.path.dirname(os.path.abspath(__file__))
cli = Client(net_profile=os.path.join(dir_path, "example_server_network.json"))
# Make the client know there is a channel in the network
cli.new_channel('mychannel')

class OperateChaincode():

    """
    Initialize common parameters: requestor, channel_name, peers, cc_name
    requestor: user role who issue the request
    channel_name: the name of the channel to send tx proposal
    peers: list of peer name and/or peer to install
    cc_name: chaincode name
    """
    def __init__(self, requestor, channel_name, peers, cc_name):
        self.requestor = requestor
        self.channel_name = channel_name
        self.peers = peers
        self.cc_name = cc_name

    #Call python SDK for query
    def query(self, args):
        return loop.run_until_complete(cli.chaincode_query(
               requestor = self.requestor,
               channel_name = self.channel_name,
               peers = self.peers,
               args = args,
               cc_name = self.cc_name
               ))

    #Call python SDK for invoke
    def invoke(self, args):
        return loop.run_until_complete(cli.chaincode_invoke(
               requestor = self.requestor,
               channel_name = self.channel_name,
               peers = self.peers,
               args = args,
               cc_name = self.cc_name,
               wait_for_event = True
               ))

    """
// =======================================================
// groupbuyRaise
// Raising(init) a groupbuy product,
// 1. Adds product to groupBuy (status=open),
// 2. Adds record to groupBuy_record,
// input: transactionID, groupbuyID, productID, currency, target_amount, dividend, expiry_date, capital
// =======================================================
  
    """
    def groupbuyRaise(self, args):
        return self.invoke(['groupbuyRaise'] + args)
    """
// =======================================================
// userRaise
// Client data raise,
// 1.Adds client to User,
// input: clientID, userID, bank, bank_account, address, phone, status_open, status_bind
// =======================================================

    """
    def userRaise(self, args):
        return self.invoke(['userRaise'] + args)

    """
// =======================================================
// groupbuyJoin
// Client joining a groupbuy product,
// 1. Update record to groupbuy.share (groupbuy.share=groupbuy.share+share)
// 2. Adds record to transcation_action (status=join),
// 3. Adds record to transcation_action_record ,
// 4. Adds record to transcation_blocked (status=blocked),
// 5. Adds record to transcation_blocked_record ,
// input: TransactionID, ClientID, ProductID, currency, share
// =======================================================

    """
    def groupbuyJoin(self,args):
        return self.invoke(['groupbuyJoin'] + args)

    """
// =======================================================
// groupbuyJoin_failed
// Client Join a groupbuy product failed,
// 1. Update record to transaction_action (status=fail),
// 2. Adds record to transaction_action_record ,
// 3. Update record to Transaction_blocked (status=fail),
// 4. Adds record to Transaction_blocked_record ,
// 5. Update record to groupbuy.share (groupbuy.share=groupbuy.share-share)
// input: TransactionID, ClientID, ProductID
// =======================================================

    """
    def groupbuyJoin_failed(self,args):
        return self.invoke(['groupbuyJoin_failed'] + args)

    """
// =======================================================
// groupbuyLeave
// Client Leaving a groupbuy product,
// 1. Adds record to transaction_action (status=leave),
// 2. Adds record to transaction_action_record ,
// 3. Update record to Transaction_blocked (status=fail),
// 4. Adds record to Transaction_blocked_record ,
// 5. Update record to groupbuy.share (groupbuy.share=groupbuy.share-share)
// input: TransactionID, ClientID, ProductID
// =======================================================

    """
    def groupbuyLeave(self,args):
        return self.invoke(['groupbuyLeave'] + args)

    """
// =======================================================
// groupbuyLeave_failed
// Client Leaved a groupbuy product fail,
// 1. Update record to transaction_action (status=join),
// 2. Adds record to transaction_action_record ,
// 3. Update record to Transaction_blocked (status=blocked),
// 4. Adds record to Transaction_blocked_record ,
// 5. Update record to groupbuy.share (groupbuy.share=groupbuy.share+share)
// input: TransactionID, ClientID, ProductID
// =======================================================

    """
    def groupbuyLeave_failed(self,args):
        return self.invoke(['groupbuyLeave_failed'] + args)

    """
// =======================================================
// groupbuyBlock
// Blocked this transactionID
// 1. Update Transaction_blocked (status = blocked),
// 2. Add to Transaction_blocked_record,
// input: transactionID, clientID, groupbuyID, currency, amount,
// =======================================================

    """
    def groupbuyBlocked(self,args):
        return self.invoke(['groupbuyBlocked'] + args)

    """
// =======================================================
// groupbuyBlocked_failed
// Blocked this transactionID failed,
// 1. Update Transaction_blocked (status=fail),
// 2. Adds record to transaction_blocked_record ,
// input: transactionID,
// =======================================================

    """
    def groupbuyBlocked_failed(self,args):
        return self.invoke(['groupbuyBlocked_failed'] + args)

    """
// =======================================================
// groupbuyUnblocked
// Unblocked this transactionID,
// 1. Update Transaction_blocked (status=Unblocked),
// 2. Adds record to transaction_blocked_record ,
// 3. Adds record toTransaction_contract (status="join_amount"),
// 4. Adds record toTransaction_contract_record,
// input: transactionID,
// =======================================================

    """
    def groupbuyUnblocked(self,args):
        return self.invoke(['groupbuyUnblocked'] + args)

    """
// =======================================================
// groupbuyUnblocked_failed
// Unblocked this transactionID failed,
// 1. Updates Transaction_blocked (status=blocked),
// 2. Adds record to transaction_blocked_record,
// 3. Updates record to Transaction_contract (status="fail"),
// 4. Adds record toTransaction_contract_record,
// input: TransactionID, ClientID, ProductID
// =======================================================

    """
    def groupbuyUnblocked_failed(self,args):
        return self.invoke(['groupbuyUnblocked_failed'] + args)

    """
// =======================================================
// groupbuySuccess
// Changing groupbuy status to success,
// 1. Updates Groupbuy (status=success),
// 2. Add record to groupbuy_record,
// input: transaction_id, groupbuyID
// =======================================================

    """
    def groupbuySuccess(self,args):
        return self.invoke(['groupbuySuccess'] + args)

    """
// =======================================================
// groupbuyNoGoodPrice
// Changing groupbuy status to no good price,
// 1. Updates Groupbuy (status=no good price),
// 2. Add record to groupbuy_record,
// input: transaction_id, groupbuyID
// =======================================================

    """
    def groupbuyNoGoodPrice(self, args):
        return self.invoke(['groupbuyNoGoodPrice'] + args)

    """
// =======================================================
// groupbuyInvesting
// Changing groupbuy status to investing,
// 1. Update Groupbuy (status=Investing),
// 2. Add record to groupbuy_record,
// input: transaction_id, groupbuyID
// =======================================================

    """
    def groupbuyInvesting(self, args):
        return self.invoke(['groupbuyInvesting'] + args)

    """
// =======================================================
// groupbuyMataned
// Changing groupbuy status to mataned,
// 1. Update Groupbuy (status=mataned),
// 2. Add record to groupbuy_record,
// input: transaction_id, groupbuyID
// =======================================================
  
    """
    def groupbuyMataned(self, args):
        return self.invoke(['groupbuyMataned'] + args)

    """
// =======================================================
// groupbuyLiquidate
// Changing groupbuy status to Liquidate,
// 1. Update Groupbuy (status=liquidate),
// 2. Add record to groupbuy_record,
// input: transaction_id, groupbuyID
// =======================================================

    """
    def groupbuyLiquidate(self, args):
        return self.invoke(['groupbuyLiquidate'] + args)

    """
// =======================================================
// groupbuyEnded
// Changing groupbuy status to ended,
// 1. Update Groupbuy (status=ended),
// 2. Add record to groupbuy_record,
// input: transaction_id, groupbuyID
// =======================================================

    """
    def groupbuyEnded(self, args):
        return self.invoke(['groupbuyEnded'] + args)

    """
// =======================================================
// groupbuyFailed
// Changing groupbuy status to fail,
// 1.Update Groupbuy (status=fail)
// 2.Add record to groupbuy_record,
// input: transaction_id, groupbuyID
// =======================================================

    """
    def groupbuyFailed(self, args):
        return self.invoke(['groupbuyFailed'] + args)

    """
// =======================================================
// dividend
// 1.Adds records to Groupbuy_contract
// 2.Adds records tp Groupbuy_contract_records
// status = dividend
// input: transactionID, clientID, groupbuyID, currency, amount,
// =======================================================

    """
    def dividend(self, args):
        return self.invoke(['dividend'] + args)

    """
// =======================================================
// Capital
// 1.Adds records to Groupbuy_contract
// 2.Adds records tp Groupbuy_contract_records
// status = capital
// input: transactionID, clientID, groupbuyID, currency, amount
// =======================================================

    """
    def capital(self, args):
        return self.invoke(['capital'] + args)

    """
// =======================================================
// liquidated
// 1.Adds records to Groupbuy_contract
// 2.Adds records tp Groupbuy_contract_records
// status = liquidate
// input: transactionID, clientID, groupbuyID, currency, amount
// =======================================================

    """
    def liquidated(self, args):
        return self.invoke(['liquidated'] + args)

    """
// =======================================================
// contract_failed
// 1.update Transaction_concract Status = fail
// 2.Add to Transaction_contract_record
// input: transactionID
// =======================================================

    """
    def contract_failed(self, args):
        return self.invoke(['contract_failed'] + args)

    """
// =======================================================
// takeover_raise
// 1.Adds records to Transaction_takeover (status=open),
// 2.Adds records to Transaction_takeover_record,
// 3.Adds records to Transaction_blocked (status=blocked),
// 4.Adds records to Transaction_blocked_record,
// input: transactionID, groupbuyID, currency, amount, clientID_sell
// =======================================================
   
    """
    def takeOver_raise(self, args):
        return self.invoke(['takeOver_raise'] + args)

    """
// =======================================================
// takeover_failed
// 1.Updates Transaction_takeover (status=fail)
// 2.Adds records to Transaction_takeover_record,
// 3.Updates Transaction_blocked (status=fail)
// 4.Add buyer transaction record to transaction
// input: ClientID_sell, ClientID_buy, ProductID.
// =======================================================

    """
    def takeOver_failed(self, args):
        return self.invoke(['takeOver_failed'] + args)
   
    """
// =======================================================
// takeover
// 1.Updates Transaction_takeover (status=success),
// 2.Adds Transaction_takeover_record,
// 3.Adds clientID_sell to Transaction_contract
// 4.Adds clientID_sell to Transaction_contract_record,
// 5.Adds clientID_buy to Transaction_contract,
// 6.Adds clientID_buy to Transaction_contract_record,
// 7.Updates Transaction_blocked (status=unblocked),
// 8.Adds Transaction_blocked_record,
// input: transactionID, groupbuyID, currency, amount, clientID_sell, clientID_buy
// =======================================================

    """
    def takeOver(self, args):
        return self.invoke(['takeOver'] + args)



    """
    queryTransaction_blockedByTransactionID

    """
    def queryTransaction_blockedByTransactionID(self, args):
        return self.invoke(['queryTransaction_blockedByTransactionID'] + args)

    """
    queryTransaction_blockedByClientID

    """
    def queryTransaction_blockedByClientID(self, args):
        return self.invoke(['queryTransaction_blockedByClientID'] + args)

    """
    queryTransaction_blocked_recordByTransactionID

    """
    def queryTransaction_blocked_recordByTransactionID(self, args):
        return self.invoke(['queryTransaction_blocked_recordByTransactionID'] + args)

    """
    queryTransaction_contractByGroupbuyID

    """
    def queryTransaction_contractByGroupbuyID(self, args):
        return self.invoke(['queryTransaction_contractByGroupbuyID'] + args)

    """
    queryTransaction_contractByClientID

    """
    def queryTransaction_contractByClientID(self, args):
        return self.invoke(['queryTransaction_contractByClientID'] + args)

    """
    queryTransaction_contractByTransactionID

    """
    def queryTransaction_contractByTransactionID(self, args):
        return self.invoke(['queryTransaction_contractByTransactionID'] + args)

    """
    queryTransaction_action_joinByClientID

    """
    def queryTransaction_action_joinByClientID(self, args):
        return self.invoke(['queryTransaction_action_joinByClientID'] + args)

    """
    queryTransaction_action_leaveByClientID

    """
    def queryTransaction_action_leaveByClientID(self, args):
        return self.invoke(['queryTransaction_action_leaveByClientID'] + args)

    """
    queryTransaction_action_joinByGroupbuyID

    """
    def queryTransaction_action_joinByGroupbuyID(self, args):
        return self.invoke(['queryTransaction_action_joinByGroupbuyID'] + args)

    """
    queryTransaction_action_leaveByGroupbuyID

    """
    def queryTransaction_action_leaveByGroupbuyID(self, args):
        return self.invoke(['queryTransaction_action_leaveByGroupbuyID'] + args)

    """
    queryTransaction_actionByTransactionID

    """
    def queryTransaction_actionByTransactionID(self, args):
        return self.invoke(['queryTransaction_actionByTransactionID'] + args)

    """
    queryTransaction_action_recordByTransactionID

    """
    def queryTransaction_action_recordByTransactionID(self, args):
        return self.invoke(['queryTransaction_action_recordByTransactionID'] + args)

    """
    queryTransaction_takeoverByClientID_sell

    """
    def queryTransaction_takeoverByClientID_sell(self, args):
        return self.invoke(['queryTransaction_takeoverByClientID_sell'] + args)

    """
    queryTransaction_takeoverByClientID_buy

    """
    def queryTransaction_takeoverByClientID_buy(self, args):
        return self.invoke(['queryTransaction_takeoverByClientID_buy'] + args)

    """
    queryTransaction_takeoverByGroupbuyID

    """
    def queryTransaction_takeoverByGroupbuyID(self, args):
        return self.invoke(['queryTransaction_takeoverByGroupbuyID'] + args)

    """
    queryTransaction_takeoverByTransactionID

    """
    def queryTransaction_takeoverByTransactionID(self, args):
        return self.invoke(['queryTransaction_takeoverByTransactionID'] + args)

    """
    queryTransaction_takeover_recordByTransactionID

    """
    def queryTransaction_takeover_recordByTransactionID(self, args):
        return self.invoke(['queryTransaction_takeover_recordByTransactionID'] + args)

    """
    queryGroupbuyByGroupbuyID

    """
    def queryGroupbuyByGroupbuyID(self, args):
        return self.invoke(['queryGroupbuyByGroupbuyID'] + args)

    """
    queryGroupbuy_recordByGroupbuyID

    """
    def queryGroupbuy_recordByGroupbuyID(self, args):
        return self.invoke(['queryGroupbuy_recordByGroupbuyID'] + args)

    """
    queryGroupbuy_recordByGroupbuyID

    """
    def queryGroupbuy_recordByTransactionID(self, args):
        return self.invoke(['queryGroupbuy_recordByTransactionID'] + args)

    """
    queryUserByClientID

    """
    def queryUserByClientID(self, args):
        return self.invoke(['queryUserByClientID'] + args)




    """
    All User.csv

    """
    def queryUser(self):
        return self.invoke(['queryUser'])

    """
    All Groupbuy.csv

    """
    def queryGroupbuy(self):
        return self.invoke(['queryGroupbuy'])

    """
    All Groupbuy_record.csv

    """
    def queryGroupbuy_record(self, args):
        return self.invoke(['queryGroupbuy_record'])

    """
    All Transaction_blocked.csv

    """
    def queryTransaction_blocked(self):
        return self.invoke(['queryTransaction_blocked'])

    """
    All Transaction_blocked_record.csv

    """
    def queryTransaction_blocked_record(self):
        return self.invoke(['queryTransaction_blocked_record'])

    """
    All Transaction_contract.csv

    """
    def queryTransaction_contract(self):
        return self.invoke(['queryTransaction_contract'])

    """
    All Transaction_contract_record.csv

    """
    def queryTransaction_contract_record(self):
        return self.invoke(['queryTransaction_contract_record'])

    """
    All Transaction_action.csv

    """
    def queryTransaction_action(self):
        return self.invoke(['queryTransaction_action'])

    """
    All Transaction_action_record.csv

    """
    def queryTransaction_action_record(self):
        return self.invoke(['queryTransaction_action_record'])

    """
    All Transaction_takeover.csv

    """
    def queryTransaction_takeover(self):
        return self.invoke(['queryTransaction_takeover'])

    """
    All Transaction_takeover_record.csv

    """
    def queryTransaction_takeover_record(self):
        return self.invoke(['queryTransaction_takeover_record'])
