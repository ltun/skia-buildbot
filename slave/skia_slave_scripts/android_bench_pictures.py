#!/usr/bin/env python
# Copyright (c) 2013 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

""" Run the Skia bench_pictures executable. """

from android_build_step import AndroidBuildStep
from build_step import BuildStep
from bench_pictures import BenchPictures
import sys


class AndroidBenchPictures(AndroidBuildStep, BenchPictures):
  def __init__(self, timeout=134400, **kwargs):
    super(AndroidBenchPictures, self).__init__(timeout=timeout, **kwargs)


if '__main__' == __name__:
  sys.exit(BuildStep.RunBuildStep(AndroidBenchPictures))
