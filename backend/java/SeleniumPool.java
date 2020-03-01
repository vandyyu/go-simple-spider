package server;
import org.openqa.selenium.Alert;
import org.openqa.selenium.By;
import org.openqa.selenium.Cookie;
import org.openqa.selenium.JavascriptExecutor;
import org.openqa.selenium.WebDriver;
import org.openqa.selenium.WebElement;
import org.openqa.selenium.firefox.FirefoxDriver;
import org.openqa.selenium.firefox.FirefoxProfile;
import org.openqa.selenium.firefox.FirefoxOptions;
import org.openqa.selenium.interactions.Actions;
import org.openqa.selenium.support.ui.Select;
import org.openqa.selenium.remote.DesiredCapabilities;

import java.util.Timer;
import java.util.TimerTask;

// NOTE: ever Selenium/FirefoxDriver consume over 150M space.
// in other words, you need 2GB when you want create 10 FirefoxDriver.
public class SeleniumPool{
	private Timer timer;
	private int cap ;
	private SeleniumExecutor[] executors;

	public SeleniumPool(int cap){
		if(cap <= 0){
			cap = 1;
		}
		this.cap = cap;
		this.executors = new SeleniumExecutor[cap];
	}
	private void statusMonitor(){
		this.timer = new Timer();
		final int cap = SeleniumPool.this.getCap();
		final int freeNum = SeleniumPool.this.getFreeNumber();
		timer.schedule(new TimerTask(){
			public void run(){
				System.out.printf("[Selenium Pool number] [used/Total] [%d/%d]\r", (cap - freeNum), cap);
			}
		}, 0, 100L);
	}

	public FirefoxOptions getFirefoxOptions(){
		System.setProperty(FirefoxDriver.SystemProperty.BROWSER_LOGFILE, "/dev/null");

		System.setProperty("webdriver.firefox.bin", "/usr/bin/firefox");
		System.setProperty("webdriver.gecko.driver", "/home/vandy/codes/go/backend/java/res/geckodriver");

		FirefoxProfile profile = new FirefoxProfile();
		// profile.setPreference("javascript.enabled", false);
		profile.setPreference("permissions.default.image", 2);
		profile.setPreference("browser.migration.version", 9001);
		profile.setPreference("permissions.default.stylesheet", 2);
		profile.setPreference("dom.ipc.plugins.enabled.libflashplayer.so", false);

		/*
		profile.setPreference("network.http.use-cache", false); 
		profile.setPreference("browser.cache.memory.enable", false);
		profile.setPreference("browser.cache.disk.enable", false);
		profile.setPreference("browser.sessionhistory.max_total_viewers", 3);
		*/
		profile.setPreference("network.dns.disableIPv6", true);   
		profile.setPreference("Content.notify.interval", 750000);
		profile.setPreference("content.notify.backoffcount", 3);
											  
		profile.setPreference("network.http.pipelining", true);                
		profile.setPreference("network.http.proxy.pipelining", true);
		profile.setPreference("network.http.pipelining.maxrequests", 32);

		FirefoxOptions options = new FirefoxOptions().setProfile(profile);
		options.addArguments("blink-settings=imagesEnabled=false");
		options.addArguments("--disable-plugins", "--disable-images","--start-maximized","--disable-javascript");
		options.addArguments("--headless");
		
		DesiredCapabilities capabilities = DesiredCapabilities.firefox();
		
		// capabilities.setCapability("pageLoadStrategy", "none");
		capabilities.setCapability("pageLoadStrategy", "eager");
		// capabilities.setCapability("pageLoadStrategy", "normal");

		options.merge(capabilities);
		return options;
	}
	public void initFirefoxPool(){
		for(int i = 0;i < this.cap;i++){
			WebDriver driver = new FirefoxDriver(this.getFirefoxOptions());
			this.executors[i] = new SeleniumExecutor(this, driver);
		}
		this.statusMonitor();
	}
	public int getFreeExecutorIndex(){
		for(int i = 0;i < this.cap;i++){
			if(this.executors[i].getTag()){
				return i;
			}
		}
		return -1;
	}
	public synchronized SeleniumExecutor getFreeExecutor(){
		for(int i = 0;i < this.cap;i++){
			if(this.executors[i].getTag()){
				this.executors[i].init();
				this.executors[i].setTag(false);
				return this.executors[i];
			}
		}
		int index = -1; 
		synchronized(this){
			try{
				// avoid fake wakeup
				while(true){
					index = this.getFreeExecutorIndex();
					if(index == -1){
						this.wait();
					}else{
						break;
					}
				}
			}catch(Exception e){
				e.printStackTrace();
			}
		}
		this.executors[index].init();
		this.executors[index].setTag(false);
		return this.executors[index];
	}
	public void freeExecutor(){
		synchronized(this){
			try{
				this.notifyAll();
			}catch(Exception e){
				e.printStackTrace();
			}
		}
	}

	public synchronized int getFreeNumber(){
		int sum = 0;
		for(int i = 0;i < this.cap;i++){
			sum += (this.executors[i].getTag() ? 1 : 0);
		}
		return sum;
	}
	public int getCap(){
		return this.cap;
	}
	public void destroy(){
		for(int i = 0;i < this.cap; i++){
			if(!this.executors[i].getTag()){
				this.executors[i].quit();
			}
		}
		if(this.timer != null){
			this.timer.cancel();
		}
	}
}
