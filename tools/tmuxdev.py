#!/usr/bin/python3

import os

import libtmux
import tmuxp


cwd = os.getcwd()
layoutstr = '5004,269x72,0,0[269x43,0,0{134x43,0,0[134x21,0,0,32,134x21,0,22,38],134x43,135,0,37},269x28,0,44{134x28,0,44,33,134x28,135,44,34}]'




def openrunner():
    ts = libtmux.server.Server()
    win = ts.sessions[0].new_window(attach=False,start_directory=cwd)
    left_pane = win.split_window(attach=False, start_directory=os.getcwd(),shell="./coolors")
    right_pane = win.split_window(attach=False, start_directory=os.getcwd(),shell="")
    win.select_layout('da26,269x72,0,0[269x43,0,0,32,269x28,0,44{134x28,0,44,33,134x28,135,44,34}]')
    
