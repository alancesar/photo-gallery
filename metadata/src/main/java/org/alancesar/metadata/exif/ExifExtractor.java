package org.alancesar.metadata.exif;

import java.io.File;
import java.util.Map;

public interface ExifExtractor {
    Map<String, String> extract(File file) throws Exception;
}
