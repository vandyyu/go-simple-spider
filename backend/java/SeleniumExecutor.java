package server;
import java.util.concurrent.TimeUnit;
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

public class SeleniumExecutor{
	private SeleniumPool pool;
	private WebDriver driver;
	// mark whether this executor is free.
	private boolean tag;
	public SeleniumExecutor(SeleniumPool pool, WebDriver driver){
		this.pool = pool;
		this.driver = driver;
		this.tag = true;
	}

	public void init(){
		((JavascriptExecutor) driver).executeScript("window.stop();");
	}
	public synchronized void setTag(boolean v){
		this.tag = v;
	}
	public synchronized boolean getTag(){
		return this.tag;
	}
	public String getHTMLData(String url){
		if (url == null || "".equals(url)){
			return null;
		}
		System.out.println("\nstart of crawling: " + url + "\n");
		String src = null;
	        try {
			this.driver.get(url);
			driver.manage().timeouts().implicitlyWait(5, TimeUnit.SECONDS);
			src = this.driver.getPageSource();
		} catch (Exception e) {
			// TODO: cannot get html data.. how to deal with it?
			e.printStackTrace();
		}finally{
			this.setTag(true);
			this.pool.freeExecutor();
		}
		System.out.println("\nend of crawling: " + url + "\n");
		return src;
	}

	public void quit(){
		if(this.driver != null){
			this.driver.quit();
		}
	}
}
