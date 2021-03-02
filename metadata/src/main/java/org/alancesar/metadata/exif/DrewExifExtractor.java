package org.alancesar.metadata.exif;

import com.drew.imaging.ImageMetadataReader;
import com.drew.imaging.ImageProcessingException;
import org.springframework.stereotype.Component;

import java.io.File;
import java.io.IOException;
import java.util.HashMap;
import java.util.Map;

@Component
public class DrewExifExtractor implements ExifExtractor {
    public Map<String, String> extract(File file) throws ImageProcessingException, IOException {
        final var tags = new HashMap<String, String>();
        final var metadata = ImageMetadataReader.readMetadata(file);
        metadata.getDirectories()
                .forEach(directory -> directory.getTags()
                        .forEach(tag -> tags.put(tag.getTagName(), tag.getDescription())));
        return tags;
    }
}
