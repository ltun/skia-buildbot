"use strict"

// Declares the testdata object that is a container for JS
// testdata used by the *-demo.html pages.

var testdata = testdata || {};

// Runs is a list of CT Pixel Diff runIDs.
testdata.runs = [
  "lchoi-20170727123456",
  "rmistry-201717202555",
  "lchoi-20170727123456",
  "rmistry-201717202555",
  "lchoi-20170727123456",
  "rmistry-201717202555",
  "lchoi-20170727123456",
  "rmistry-201717202555",
  "lchoi-20170727123456",
  "rmistry-201717202555",
  "lchoi-20170727123456",
  "rmistry-201717202555",
  "lchoi-20170727123456",
  "rmistry-201717202555",
  "lchoi-20170727123456",
  "rmistry-201717202555",
  "lchoi-20170727123456",
  "rmistry-201717202555",
  "lchoi-20170727123456",
  "rmistry-201717202555",
  "lchoi-20170727123456",
  "rmistry-201717202555",
  "lchoi-20170727123456",
  "rmistry-201717202555",
  "lchoi-20170727123456",
  "rmistry-201717202555",
  "lchoi-20170727123456",
  "rmistry-201717202555",
  "lchoi-20170727123456",
  "rmistry-201717202555",
  "lchoi-20170727123456",
  "rmistry-201717202555",
  "lchoi-20170727123456",
  "rmistry-201717202555",
  "lchoi-20170727123456",
  "rmistry-201717202555",
]

// Results is a serialized list of ResultRecs.
testdata.results = [
  {
    "URL": "http://www.google.com",
    "Rank": 1,
    "NoPatchImg": "lchoi-20170727123456/nopatch/1/http___www_google_com",
    "WithPatchImg": "lchoi-20170727123456/withpatch/1/http___www_google_com",
    "DiffMetrics": {
      "numDiffPixels": 1000,
      "pixelDiffPercent": 50.2,
      "maxRGBADiffs": [0, 128, 255, 0]
    }
  },
  {
    "URL": "http://www.youtube.com",
    "Rank": 2,
    "NoPatchImg": "lchoi-20170727123456/nopatch/2/http___www_youtube_com",
    "WithPatchImg": "lchoi-20170727123456/withpatch/2/http___www_youtube_com",
    "DiffMetrics": {
      "numDiffPixels": 2500,
      "pixelDiffPercent": 25.7,
      "maxRGBADiffs": [7, 14, 21, 0]
    }
  },
  {
    "URL": "http://www.google.com",
    "Rank": 1,
    "NoPatchImg": "lchoi-20170727123456/nopatch/1/http___www_google_com",
    "WithPatchImg": "lchoi-20170727123456/withpatch/1/http___www_google_com",
    "DiffMetrics": {
      "numDiffPixels": 1000,
      "pixelDiffPercent": 50.2,
      "maxRGBADiffs": [0, 128, 255, 0]
    }
  },
  {
    "URL": "http://www.youtube.com",
    "Rank": 2,
    "NoPatchImg": "lchoi-20170727123456/nopatch/2/http___www_youtube_com",
    "WithPatchImg": "lchoi-20170727123456/withpatch/2/http___www_youtube_com",
    "DiffMetrics": {
      "numDiffPixels": 2500,
      "pixelDiffPercent": 25.7,
      "maxRGBADiffs": [7, 14, 21, 0]
    }
  },
  {
    "URL": "http://www.google.com",
    "Rank": 1,
    "NoPatchImg": "lchoi-20170727123456/nopatch/1/http___www_google_com",
    "WithPatchImg": "lchoi-20170727123456/withpatch/1/http___www_google_com",
    "DiffMetrics": {
      "numDiffPixels": 1000,
      "pixelDiffPercent": 50.2,
      "maxRGBADiffs": [0, 128, 255, 0]
    }
  },
  {
    "URL": "http://www.youtube.com",
    "Rank": 2,
    "NoPatchImg": "lchoi-20170727123456/nopatch/2/http___www_youtube_com",
    "WithPatchImg": "lchoi-20170727123456/withpatch/2/http___www_youtube_com",
    "DiffMetrics": {
      "numDiffPixels": 2500,
      "pixelDiffPercent": 25.7,
      "maxRGBADiffs": [7, 14, 21, 0]
    }
  },
  {
    "URL": "http://www.google.com",
    "Rank": 1,
    "NoPatchImg": "lchoi-20170727123456/nopatch/1/http___www_google_com",
    "WithPatchImg": "lchoi-20170727123456/withpatch/1/http___www_google_com",
    "DiffMetrics": {
      "numDiffPixels": 1000,
      "pixelDiffPercent": 50.2,
      "maxRGBADiffs": [0, 128, 255, 0]
    }
  },
  {
    "URL": "http://www.youtube.com",
    "Rank": 2,
    "NoPatchImg": "lchoi-20170727123456/nopatch/2/http___www_youtube_com",
    "WithPatchImg": "lchoi-20170727123456/withpatch/2/http___www_youtube_com",
    "DiffMetrics": {
      "numDiffPixels": 2500,
      "pixelDiffPercent": 25.7,
      "maxRGBADiffs": [7, 14, 21, 0]
    }
  }
]
