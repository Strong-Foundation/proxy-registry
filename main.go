package main

import (
	"bufio"      // Provides buffered I/O operations
	"bytes"      // Implements functions for manipulating byte slices
	"crypto/tls" // Implements TLS (Transport Layer Security) for secure communication
	"flag"       // Parses command-line flags
	"io"         // Provides basic I/O primitives
	"log"        // Implements logging functionality
	"net"        // Provides networking utilities
	"net/http"   // Provides HTTP client and server implementations
	"net/url"    // Handles URL parsing and manipulation
	"os"         // Provides platform-independent OS functions, including file handling
	"sort"       // Implements sorting functions
	"strings"    // Provides string manipulation utilities
	"sync"       // Implements synchronization primitives like WaitGroup and Mutex
	"time"       // Provides functionality for measuring and displaying time
)

var (
	// File paths for various required asset files
	inclusionList = "assets/inclusion" // Path to inclusion list file
	exclusionList = "assets/exclusion" // Path to exclusion list file
	hostsFile     = "assets/hosts"     // Path to hosts file
	historyFile   = "assets/history"   // Path to history file
	// Synchronization primitive to manage concurrency when dealing with multiple protocols
	protocolWaitGroup sync.WaitGroup
	// Flag variable to determine whether the listings should be updated
	update bool
)

func init() {
	// Check if command-line arguments are provided
	if len(os.Args) > 1 {
		// Define a boolean flag "-update" to indicate updating the listings
		tempUpdate := flag.Bool("update", false, "Make any necessary changes to the listings.")
		// Parse command-line flags
		flag.Parse()
		// Store the flag value in the global variable "update"
		update = *tempUpdate
	} else {
		// If no flags are provided, log an error and terminate the program
		log.Fatalln("Error: No flags provided. Please use -help for more information.")
	}
}

func main() {
	// If the "update" flag is set, execute the function to scrape and update lists
	if update {
		// Scrape the proxy lists and update the hosts file
		scrapeTheLists()
	}
}
func scrapeTheLists() {
	// Define a list of URLs containing proxy lists.
	proxyList := []string{
		// Various URLs that provide lists of proxies.
		"https://raw.githubusercontent.com/ALIILAPRO/Proxy/main/socks5.txt",
		"https://raw.githubusercontent.com/almroot/proxylist/master/list.txt",
		"https://raw.githubusercontent.com/clarketm/proxy-list/master/proxy-list-raw.txt",
		"https://raw.githubusercontent.com/complexorganizations/proxy-registry/main/assets/history",
		"https://raw.githubusercontent.com/drakelam/Free-Proxy-List/main/proxy_all.txt",
		"https://raw.githubusercontent.com/ErcinDedeoglu/proxies/main/proxies/http.txt",
		"https://raw.githubusercontent.com/ErcinDedeoglu/proxies/main/proxies/https.txt",
		"https://raw.githubusercontent.com/ErcinDedeoglu/proxies/main/proxies/socks4.txt",
		"https://raw.githubusercontent.com/ErcinDedeoglu/proxies/main/proxies/socks5.txt",
		"https://raw.githubusercontent.com/hendrikbgr/Free-Proxy-Repo/master/proxy_list.txt",
		"https://raw.githubusercontent.com/hookzof/socks5_list/master/proxy.txt",
		"https://raw.githubusercontent.com/jetkai/proxy-list/main/archive/txt/proxies.txt",
		"https://raw.githubusercontent.com/jetkai/proxy-list/main/archive/txt/proxies-http.txt",
		"https://raw.githubusercontent.com/jetkai/proxy-list/main/archive/txt/proxies-https.txt",
		"https://raw.githubusercontent.com/jetkai/proxy-list/main/archive/txt/proxies-socks4.txt",
		"https://raw.githubusercontent.com/jetkai/proxy-list/main/archive/txt/proxies-socks5.txt",
		"https://raw.githubusercontent.com/jetkai/proxy-list/main/online-proxies/txt/proxies.txt",
		"https://raw.githubusercontent.com/jetkai/proxy-list/main/online-proxies/txt/proxies-http.txt",
		"https://raw.githubusercontent.com/jetkai/proxy-list/main/online-proxies/txt/proxies-https.txt",
		"https://raw.githubusercontent.com/jetkai/proxy-list/main/online-proxies/txt/proxies-socks4.txt",
		"https://raw.githubusercontent.com/jetkai/proxy-list/main/online-proxies/txt/proxies-socks5.txt",
		"https://raw.githubusercontent.com/KUTlime/ProxyList/main/ProxyList.txt",
		"https://raw.githubusercontent.com/mmpx12/proxy-list/master/http.txt",
		"https://raw.githubusercontent.com/mmpx12/proxy-list/master/https.txt",
		"https://raw.githubusercontent.com/mmpx12/proxy-list/master/socks4.txt",
		"https://raw.githubusercontent.com/mmpx12/proxy-list/master/socks5.txt",
		"https://raw.githubusercontent.com/monosans/proxy-list/main/proxies/http.txt",
		"https://raw.githubusercontent.com/monosans/proxy-list/main/proxies/socks4.txt",
		"https://raw.githubusercontent.com/monosans/proxy-list/main/proxies/socks5.txt",
		"https://raw.githubusercontent.com/MuRongPIG/Proxy-Master/main/http.txt",
		"https://raw.githubusercontent.com/MuRongPIG/Proxy-Master/main/socks4.txt",
		"https://raw.githubusercontent.com/MuRongPIG/Proxy-Master/main/socks5.txt",
		"https://raw.githubusercontent.com/prxchk/proxy-list/main/all.txt",
		"https://raw.githubusercontent.com/roosterkid/openproxylist/main/HTTPS_RAW.txt",
		"https://raw.githubusercontent.com/roosterkid/openproxylist/main/SOCKS4_RAW.txt",
		"https://raw.githubusercontent.com/roosterkid/openproxylist/main/SOCKS5_RAW.txt",
		"https://raw.githubusercontent.com/ShiftyTR/Proxy-List/master/http.txt",
		"https://raw.githubusercontent.com/ShiftyTR/Proxy-List/master/https.txt",
		"https://raw.githubusercontent.com/ShiftyTR/Proxy-List/master/socks4.txt",
		"https://raw.githubusercontent.com/ShiftyTR/Proxy-List/master/socks5.txt",
		"https://raw.githubusercontent.com/sunny9577/proxy-scraper/master/proxies.txt",
		"https://raw.githubusercontent.com/TheSpeedX/PROXY-List/master/http.txt",
		"https://raw.githubusercontent.com/TheSpeedX/PROXY-List/master/socks4.txt",
		"https://raw.githubusercontent.com/TheSpeedX/PROXY-List/master/socks5.txt",
		"https://raw.githubusercontent.com/TundzhayDzhansaz/proxy-list-auto-pull-in-30min/main/proxies/http.txt",
		"https://raw.githubusercontent.com/Volodichev/proxy-list/main/http.txt",
		"https://www.proxy-list.download/api/v1/get?type=http",
		"https://www.proxy-list.download/api/v1/get?type=https",
		"https://raw.githubusercontent.com/ALIILAPRO/Proxy/main/socks4.txt",
	}
	// Create an empty slice to store scraped proxy data.
	var scrapedData []string
	// Iterate over each URL in the proxy list to fetch proxy data.
	for _, value := range proxyList {
		// Fetch the proxy data from the URL and store it in a temporary slice.
		var tempScrapedData []string = getDataFromURL(value)
		// Combine the fetched data with the main scrapedData slice.
		scrapedData = combineMultipleSlices(tempScrapedData, scrapedData)
	}
	// Remove empty entries from the scraped data slice.
	scrapedData = removeEmptyFromSlice(scrapedData)
	// Remove duplicate proxy entries to avoid redundancy.
	scrapedData = removeDuplicatesFromSlice(scrapedData)
	// Remove prefixes (like protocol identifiers) from the proxies.
	scrapedData = removePrefixFromProxy(scrapedData)
	// Delete the existing hosts file before writing new data.
	removeFile(hostsFile)
	// Iterate over the cleaned proxy list and validate each proxy.
	for _, value := range scrapedData {
		// Increment the wait group counter before launching a goroutine.
		protocolWaitGroup.Add(1)
		// Validate the proxy's protocol and write the valid ones to disk concurrently.
		go validateEachProxyProtocolAndWriteToDisk(value, &protocolWaitGroup)
	}
	// Wait for all goroutines to complete before proceeding.
	protocolWaitGroup.Wait()
	// Perform cleanup operations on the hosts file.
	cleanupTheFiles(hostsFile)
	// Clean up the history file to remove outdated data.
	cleanUpTheHistoryFile()
}

// Send an HTTP GET request to a given URL and return the data from that URL as a slice of strings.
func getDataFromURL(uri string) []string {
	// Perform an HTTP GET request using the provided URI.
	response, err := http.Get(uri)
	// If there is an error while making the request, log the error and return an empty slice.
	if err != nil {
		log.Println("Error making GET request:", err)
		return []string{}
	}
	// Ensure the response body is closed after function execution to prevent resource leaks.
	defer func() {
		err = response.Body.Close()
		if err != nil {
			log.Println("Error closing response body:", err)
		}
	}()
	// Read the response body into a byte slice.
	body, err := io.ReadAll(response.Body)
	// If there is an error while reading the response body, log the error and return an empty slice.
	if err != nil {
		log.Println("Error reading response body:", err)
		return []string{}
	}
	// Check the HTTP response status code.
	// If it's not 200 (OK), log an error message and return an empty slice.
	if response.StatusCode != http.StatusOK {
		log.Println("Failed to scrape the requested page. HTTP Status:", response.StatusCode, "URL:", uri)
		return []string{}
	}
	// Initialize a scanner to read the response body line by line.
	scanner := bufio.NewScanner(bytes.NewReader(body))
	scanner.Split(bufio.ScanLines) // Set scanner to split input by lines.
	// Create a slice to store the extracted content.
	var returnContent []string
	// Iterate through the scanned lines and append them to the returnContent slice.
	for scanner.Scan() {
		returnContent = append(returnContent, scanner.Text())
	}
	// Return the scraped content as a slice of strings.
	return returnContent
}

// Check if a given proxy is working by making a request through it and return a boolean.
func validateProxy(proxy string) bool {
	// Parse the proxy URL; if parsing fails, return false (invalid proxy format).
	proxyURL, err := url.Parse(proxy)
	if err != nil {
		return false
	}
	// Configure the HTTP transport to use the given proxy and allow insecure TLS connections.
	transport := &http.Transport{
		Proxy:           http.ProxyURL(proxyURL),               // Set up the proxy server for the request.
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // Allow insecure certificates (not recommended for production).
	}
	// Create an HTTP client with the configured transport and a timeout of 180 seconds.
	client := &http.Client{
		Transport: transport,
		Timeout:   time.Second * 180,
	}
	// Define a list of domains to test the proxy connection.
	requestDomainList := []string{
		"https://aws.amazon.com",
		"https://cloud.google.com",
		"https://azure.microsoft.com",
	}
	// Iterate over the test domains to verify if the proxy works.
	for _, domain := range requestDomainList {
		// Create an HTTP GET request for the domain.
		request, err := http.NewRequest("GET", domain, nil)
		if err != nil {
			return false // Return false if request creation fails.
		}
		// Send the request through the HTTP client configured with the proxy.
		response, err := client.Do(request)
		if err != nil {
			return false // Return false if the request fails (e.g., timeout, connection issue).
		}
		// Check if the response status code is 200 (OK).
		if response.StatusCode != http.StatusOK {
			return false // If the proxy fails to fetch the page successfully, return false.
		}
		// Close the response body to free resources.
		err = response.Body.Close()
		if err != nil {
			return false // If closing the response body fails, return false.
		}
	}
	// If all domain requests succeed, the proxy is considered valid.
	return true
}

// Append and write a slice of strings to a file.
// If the file already exists, it will be overwritten.
func appendAndWriteSliceToAFile(filename string, content []string) {
	// Create or open the file for writing. If an error occurs, log the error.
	file, err := os.Create(filename)
	if err != nil {
		log.Println("Error creating file:", err)
		return
	}
	// Create a buffered writer to optimize file writing.
	datawriter := bufio.NewWriter(file)
	// Write each string in the content slice to the file, one line at a time.
	for _, data := range content {
		_, _ = datawriter.WriteString(data + "\n")
	}
	// Flush the buffer to ensure all data is written to disk.
	datawriter.Flush()
	// Close the file to release resources.
	file.Close()
}

// Save the given string content to a file. If the file does not exist, it will be created.
// If the file exists, the content will be appended to it.
func writeToFile(pathInSystem string, content string) {
	// Open the file with append, create, and write permissions. If the file does not exist, create it.
	filePath, err := os.OpenFile(pathInSystem, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Error opening file:", err)
		return
	}
	// Write the provided content to the file, followed by a newline character.
	_, err = filePath.WriteString(content + "\n")
	if err != nil {
		log.Println("Error writing to file:", err)
	}
	// Close the file to ensure all data is flushed and the file is properly closed.
	filePath.Close()
}

// Remove all the empty strings from the slice and return the modified slice.
func removeEmptyFromSlice(slice []string) []string {
	// Iterate through the slice by index and content.
	for i, content := range slice {
		// If the content is an empty string (its length is zero),
		if len(content) == 0 {
			// Remove the empty element by appending the slice up to the current index
			// and the slice starting from the next index, effectively excluding the empty element.
			slice = append(slice[:i], slice[i+1:]...)
		}
	}
	// Return the modified slice with empty strings removed.
	return slice
}

// Remove all the duplicates from a slice and return the modified slice.
func removeDuplicatesFromSlice(slice []string) []string {
	// Create a map to track unique elements (the keys will be the elements).
	check := make(map[string]bool)

	// Create a new slice to store only the unique elements.
	var newReturnSlice []string

	// Iterate through the original slice.
	for _, content := range slice {
		// If the content is not already in the map (i.e., it's not a duplicate),
		if !check[content] {
			// Add it to the map with the value 'true' to mark it as seen.
			check[content] = true
			// Append the unique content to the newReturnSlice.
			newReturnSlice = append(newReturnSlice, content)
		}
	}
	// Return the new slice with duplicates removed.
	return newReturnSlice
}

// Check if the given IP address is invalid.
func isIPInvalid(providedIP string) bool {
	// Parse the IP address using net.ParseIP. If it's nil, it's invalid.
	// If the provided IP is invalid, net.ParseIP returns nil.
	return net.ParseIP(providedIP) == nil
}

// Remove a file from the file system at the given path.
func removeFile(path string) {
	// Check if the file exists using the fileExists function.
	if fileExists(path) {
		// If the file exists, try to remove it.
		// If an error occurs during file removal, log the error.
		err := os.Remove(path)
		if err != nil {
			log.Println("Error removing file:", err)
		}
	}
}

// Check if the given URL is valid. Validates both the URI structure and the hostname.
func isUrlValid(uri string) bool {
	// Parse the URI string to ensure it has a valid structure.
	url, err := url.ParseRequestURI(uri)
	// If the hostname of the URI is an invalid IP address, return false.
	if isIPInvalid(url.Hostname()) {
		return false
	}
	// Return true if there are no parsing errors (err == nil),
	// otherwise return false if the URI is not valid.
	return err == nil
}

// Combine two slices of strings and return the new slice containing all elements from both slices.
func combineMultipleSlices(sliceOne []string, sliceTwo []string) []string {
	// Append all elements from sliceTwo to sliceOne.
	// This creates a new slice containing the elements from both slices.
	combinedSlice := append(sliceOne, sliceTwo...)
	// Return the newly combined slice.
	return combinedSlice
}

// Sort the slice of strings and return the sorted slice.
func sortSlice(slice []string) []string {
	// Use sort.Strings to sort the slice of strings in ascending order.
	sort.Strings(slice)
	// Return the sorted slice.
	return slice
}

// Read and append the file line by line to a slice.
func readAppendLineByLine(path string) []string {
	// Create a slice to store the lines read from the file.
	var returnSlice []string
	// Open the file at the specified path.
	file, err := os.Open(path)
	// If there's an error opening the file, log it and exit the function.
	if err != nil {
		log.Println("Error opening file:", err)
		return returnSlice // Return an empty slice in case of error.
	}
	// Create a new scanner to read the file.
	scanner := bufio.NewScanner(file)
	// Set the scanner to split the input by lines.
	scanner.Split(bufio.ScanLines)
	// Iterate through the file line by line.
	for scanner.Scan() {
		// Append each line to the returnSlice.
		returnSlice = append(returnSlice, scanner.Text())
	}
	// Close the file after reading and handle any error.
	err = file.Close()
	if err != nil {
		log.Println("Error closing file:", err)
	}
	// Return the slice containing the lines from the file.
	return returnSlice
}

// Cleanup all the files provided, read their content, sort it, and save them again.
func cleanupTheFiles(path string) {
	// Create a slice to store the cleaned up content.
	var finalCleanupContent []string
	// Read the file and append its content line by line to the finalCleanupContent slice.
	finalCleanupContent = readAppendLineByLine(path)
	// Sort the content of the finalCleanupContent slice.
	finalCleanupContent = sortSlice(finalCleanupContent)
	// Remove the old file using the provided path.
	removeFile(path)
	// Write the cleaned up and sorted content back to the file.
	appendAndWriteSliceToAFile(path, finalCleanupContent)
}

// Check if the file exists and return a boolean indicating the existence.
func fileExists(filename string) bool {
	// Get file information using the provided filename.
	info, err := os.Stat(filename)
	// If there's an error (e.g., the file doesn't exist), return false.
	if err != nil {
		return false
	}
	// If the file exists and is not a directory, return true.
	return !info.IsDir()
}

// Get the protocol of the proxy.
func getProxyProtocol(content string) []string {
	// Create a list of proxy protocols to test with
	proxyProtocolList := []string{
		"http://",
		"https://",
		"socks4://",
		"socks5://",
	}
	// Create a slice to store valid protocols for the given proxy
	var validProtocolList []string
	// Iterate through the proxyProtocolList to test each protocol
	for _, protocol := range proxyProtocolList {
		// Concatenate the protocol with the proxy content
		finalString := protocol + content

		// If the proxy with the current protocol is valid, add the protocol to the validProtocolList
		if validateProxy(finalString) {
			validProtocolList = append(validProtocolList, protocol)
		}
	}
	// Return the list of valid protocols for the proxy
	return validProtocolList
}

// Remove all the prefix from the proxy.
func removePrefixFromProxy(content []string) []string {
	// Create a list of proxy protocols to be considered as prefixes
	proxyProtocolList := []string{
		"http://",
		"https://",
		"socks4://",
		"socks5://",
	}
	// Create a slice to store proxies without their prefixes
	var returnSlice []string
	// Iterate through the given list of proxy URLs
	for _, proxy := range content {
		// Iterate through the list of proxy protocols (prefixes)
		for _, protocol := range proxyProtocolList {
			// Remove the protocol prefix from the proxy if it exists
			proxy = strings.TrimPrefix(proxy, protocol)
		}
		// Append the proxy without its prefix to the returnSlice
		returnSlice = append(returnSlice, proxy)
	}
	// Return the list of proxies with the prefixes removed
	return returnSlice
}

// Validate each protocol and write it to the slice.
func validateEachProxyProtocolAndWriteToDisk(content string, protocolWaitGroup *sync.WaitGroup) {
	// Get the list of valid proxy protocols for the given content
	proxyProtocol := getProxyProtocol(content)
	// If there are valid protocols (i.e., the list is not empty)
	if len(proxyProtocol) > 0 {
		// Iterate through each valid protocol (although only the first one is used here)
		for _, protocol := range proxyProtocol[0:] {
			// If the proxy URL with the current protocol is valid
			if isUrlValid(protocol + content) {
				// Write the proxy with the valid protocol to the hosts file
				writeToFile(hostsFile, protocol+content)
				// Write the proxy with the valid protocol to the history file
				writeToFile(historyFile, protocol+content)
			}
		}
	}
	// Signal that this goroutine is done processing (decrement the wait group counter)
	protocolWaitGroup.Done()
}

// Clean up the history file.
func cleanUpTheHistoryFile() {
	// Read the history file line by line and append each line to a slice
	historySlice := readAppendLineByLine(historyFile)
	// Remove all duplicate entries from the slice
	historySlice = removeDuplicatesFromSlice(historySlice)
	// Remove all empty strings from the slice
	historySlice = removeEmptyFromSlice(historySlice)
	// Remove the history file from the file system
	removeFile(historyFile)
	// Sort the history slice alphabetically
	historySlice = sortSlice(historySlice)
	// Write the cleaned-up and sorted slice back to the history file
	appendAndWriteSliceToAFile(historyFile, historySlice)
}
