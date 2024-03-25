# ryu_controller.py
from ryu.base import app_manager
from ryu.controller import ofp_event
from ryu.controller.handler import CONFIG_DISPATCHER, MAIN_DISPATCHER
from ryu.controller.handler import set_ev_cls
from ryu.ofproto import ofproto_v1_3

class MyController(app_manager.RyuApp):
    OFP_VERSIONS = [ofproto_v1_3.OFP_VERSION]

    def __init__(self, *args, **kwargs):
        super(MyController, self).__init__(*args, **kwargs)

    @set_ev_cls(ofp_event.EventOFPSwitchFeatures, CONFIG_DISPATCHER)
    def switch_features_handler(self, ev):
        datapath = ev.msg.datapath
        ofproto = datapath.ofproto
        parser = datapath.ofproto_parser

        # 设置规则
        # node-1与node-4
        self.add_flow(datapath, 10, parser.OFPMatch(in_port=2), [parser.OFPActionOutput(3)])
        self.add_flow(datapath, 10, parser.OFPMatch(in_port=3), [parser.OFPActionOutput(2)])

        # node-2与node-3 
        self.add_flow(datapath, 10, parser.OFPMatch(in_port=4), [parser.OFPActionOutput(8)])
        self.add_flow(datapath, 10, parser.OFPMatch(in_port=8), [parser.OFPActionOutput(4)])

        # 默认规则，阻止所有其他流量
        match = parser.OFPMatch()
        actions = []
        self.add_flow(datapath, 0, match, actions)

    def add_flow(self, datapath, priority, match, actions):
        ofproto = datapath.ofproto
        parser = datapath.ofproto_parser

        inst = [parser.OFPInstructionActions(ofproto.OFPIT_APPLY_ACTIONS, actions)]
        mod = parser.OFPFlowMod(
            datapath=datapath, priority=priority,
            match=match, instructions=inst)
        datapath.send_msg(mod)
