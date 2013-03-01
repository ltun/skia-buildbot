#!/usr/bin/env python
# Copyright (c) 2013 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

""" Download the SKPs. """

from build_step import BuildStep
from utils import gs_utils
from utils import sync_bucket_subdir
import compare_and_upload_webpage_gms
import os
import sys


class DownloadSKPs(BuildStep):
  def _CreateLocalStorageDirs(self):
    """Creates required local storage directories for this script."""
    if not os.path.exists(self._local_playback_dirs.PlaybackSkpDir()):
      os.makedirs(self._local_playback_dirs.PlaybackSkpDir())

    if os.path.exists(self._local_playback_dirs.PlaybackGmActualDir()):
      # Delete everything except the timestamp and last comparison files.
      for path, _dirs, files in os.walk(
          self._local_playback_dirs.PlaybackGmActualDir()):
        if gs_utils.TIMESTAMP_COMPLETED_FILENAME in files:
          files.remove(gs_utils.TIMESTAMP_COMPLETED_FILENAME)
        if compare_and_upload_webpage_gms.LAST_COMPARISON_FILENAME in files:
          files.remove(compare_and_upload_webpage_gms.LAST_COMPARISON_FILENAME)
        for gm_actual_file in files:
          os.remove(os.path.join(path, gm_actual_file))
    else:
      os.makedirs(self._local_playback_dirs.PlaybackGmActualDir())

  def _DownloadSKPsFromStorage(self):
    """Copies over skp files from Google Storage if the timestamps differ."""
    dest_gsbase = (self._args.get('dest_gsbase') or
                   sync_bucket_subdir.DEFAULT_PERFDATA_GS_BASE)
    print '\n\n========Downloading skp files from Google Storage========\n\n'
    gs_utils.DownloadDirectoryContentsIfChanged(
        gs_base=dest_gsbase,
        gs_relative_dir=self._storage_playback_dirs.PlaybackSkpDir(),
        local_dir=self._local_playback_dirs.PlaybackSkpDir())

  def _Run(self):
    if not self._use_skp_playback_framework:
      return

    # Create the required local storage directories.
    self._CreateLocalStorageDirs()

    # Locally copy skps generated by webpages_playback from GoogleStorage.
    self._DownloadSKPsFromStorage()


if '__main__' == __name__:
  sys.exit(BuildStep.RunBuildStep(DownloadSKPs))