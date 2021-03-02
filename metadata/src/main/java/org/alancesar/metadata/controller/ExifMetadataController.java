package org.alancesar.metadata.controller;

import org.alancesar.metadata.exif.ExifExtractor;
import org.alancesar.metadata.service.storage.StorageService;
import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.*;

import java.io.FileOutputStream;
import java.nio.file.Files;
import java.util.Map;

@RestController
@RequestMapping("exif")
public class ExifMetadataController {
    public static final String TEMP_PREFIX = "exif-";
    public static final String TEMP_SUFFIX = ".tmp";
    public static final Map<String, String> ERROR_MESSAGE = Map.of("message", "file not found");

    private final StorageService service;
    private final ExifExtractor extractor;

    public ExifMetadataController(StorageService service, ExifExtractor extractor) {
        this.service = service;
        this.extractor = extractor;
    }

    @GetMapping("{filename}")
    @ResponseBody
    public Map<String, String> getExif(@PathVariable String filename) throws Exception {
        var inputStream = service.get(filename);
        var tempFile = Files.createTempFile(TEMP_PREFIX, TEMP_SUFFIX);
        var outputStream = new FileOutputStream(tempFile.toFile());
        inputStream.transferTo(outputStream);
        return extractor.extract(tempFile.toFile());
    }

    @ExceptionHandler(Exception.class)
    @ResponseStatus(HttpStatus.NOT_FOUND)
    public Map<String, String> exceptionHandler() {
        return ERROR_MESSAGE;
    }
}
